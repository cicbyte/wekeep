package parser

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"golang.org/x/net/html"
)

// WechatArticle 解析后的微信文章结构
type WechatArticle struct {
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Content     string    `json:"content"`
	BaseURL     string    `json:"baseUrl"`
	PublishTime time.Time `json:"publishTime"`
}

// ParseWechatArticle 解析微信公众号文章
func ParseWechatArticle(htmlContent, baseURL string) (*WechatArticle, error) {
	// 替换 &nbsp; 字符
	htmlContent = strings.ReplaceAll(htmlContent, "&nbsp;", " ")

	// 解析 HTML
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("解析 HTML 失败: %v", err)
	}

	// 提取标题
	title := extractTitle(doc)

	// 提取作者
	author := extractAuthor(doc)

	// 提取发布时间
	publishTime := extractPublishTime(doc)

	// 提取文章内容（HTML 格式）
	contentHTML, err := extractContent(doc)
	if err != nil {
		return nil, err
	}

	// 处理 HTML 内容
	contentHTML = processContent(contentHTML)

	// 将 HTML 转换为 Markdown
	content := htmlToMarkdown(contentHTML)

	return &WechatArticle{
		Title:       title,
		Author:      author,
		Content:     content,
		BaseURL:     baseURL,
		PublishTime: publishTime,
	}, nil
}

// htmlToMarkdown 将 HTML 内容转换为 Markdown 格式
func htmlToMarkdown(htmlContent string) string {
	// 使用 html-to-markdown 库进行转换
	markdown, err := htmltomarkdown.ConvertString(htmlContent)
	if err != nil {
		// 如果转换失败，返回原始内容
		return htmlContent
	}

	// 清理多余的空行
	return cleanMarkdownWhitespace(markdown)
}

// cleanMarkdownWhitespace 清理 Markdown 中的多余空白
func cleanMarkdownWhitespace(content string) string {
	// 移除多余的连续空行（超过2个换行符）
	for strings.Contains(content, "\n\n\n") {
		content = strings.ReplaceAll(content, "\n\n\n", "\n\n")
	}

	// 移除首尾空白
	content = strings.TrimSpace(content)

	return content
}

// 提取页面标题（优化版）
func extractTitle(doc *html.Node) string {
	var title string

	// 按优先级提取标题：
	// 1. meta标签（最准确）
	// 2. 微信专用标题标签
	// 3. title标签

	// 方法1：从meta标签获取（优先级最高）
	metaSelectors := []map[string]string{
		{"property": "og:title"},
		{"name": "og:title"},
		{"property": "title"},
		{"name": "title"},
	}

	for _, selector := range metaSelectors {
		title = findMetaContent(doc, selector)
		if title != "" {
			return strings.TrimSpace(title)
		}
	}

	// 方法2：从微信专用标签获取
	wxSelectors := []struct {
		tag  string
		attr string
		val  string
	}{
		{"h1", "class", "rich_media_title"},
		{"h1", "id", "activity-name"},
		{"div", "class", "rich_media_title"},
	}

	for _, selector := range wxSelectors {
		title = findElementText(doc, selector.tag, selector.attr, selector.val)
		if title != "" {
			return strings.TrimSpace(title)
		}
	}

	// 方法3：从title标签获取
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			if n.FirstChild != nil {
				title = strings.TrimSpace(n.FirstChild.Data)
			}
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)

	// 去除标题中的 " - 微信" 后缀
	if strings.Contains(title, " - 微信") {
		title = strings.Split(title, " - 微信")[0]
	}

	return title
}

// 提取作者
func extractAuthor(doc *html.Node) string {
	// 微信公众号作者位置：
	// 1. meta标签
	// 2. #post-user
	// 3. .rich_media_meta_nickname

	// 尝试从 meta 标签
	if author := findMetaContent(doc, map[string]string{"property": "og:article:author"}); author != "" {
		return author
	}
	if author := findMetaContent(doc, map[string]string{"name": "og:article:author"}); author != "" {
		return author
	}

	// 从用户信息链接
	if author := findElementText(doc, "a", "id", "post-user"); author != "" {
		return author
	}

	// 从昵称标签
	if author := findElementTextByClass(doc, "span", "rich_media_meta_nickname"); author != "" {
		return author
	}

	return ""
}

// 提取发布时间
func extractPublishTime(doc *html.Node) time.Time {
	// 尝试从 meta 标签获取ISO格式时间
	if timeStr := findMetaContent(doc, map[string]string{"property": "article:published_time"}); timeStr != "" {
		if t, err := time.Parse(time.RFC3339, timeStr); err == nil {
			return t
		}
	}

	// 从发布日期标签
	timeSelectors := []struct {
		tag  string
		attr string
		val  string
	}{
		{"em", "id", "post-date"},
		{"em", "class", "rich_media_meta_text"},
	}

	for _, selector := range timeSelectors {
		if timeText := findElementText(doc, selector.tag, selector.attr, selector.val); timeText != "" {
			// 解析微信时间格式：如 "2024年1月15日"
			if t := parseWeixinTime(timeText); !t.IsZero() {
				return t
			}
		}
	}

	return time.Time{}
}

// parseWeixinTime 解析微信时间格式
func parseWeixinTime(timeStr string) time.Time {
	// 移除多余空格
	timeStr = strings.TrimSpace(timeStr)

	// 格式1：2024年1月15日
	re1 := regexp.MustCompile(`(\d{4})年(\d{1,2})月(\d{1,2})日`)
	if matches := re1.FindStringSubmatch(timeStr); matches != nil {
		if t, err := time.Parse("2006年1月2日", matches[0]); err == nil {
			return t
		}
	}

	// 格式2：2024-01-15
	re2 := regexp.MustCompile(`(\d{4})-(\d{1,2})-(\d{1,2})`)
	if matches := re2.FindStringSubmatch(timeStr); matches != nil {
		if t, err := time.Parse("2006-01-02", matches[0]); err == nil {
			return t
		}
	}

	return time.Time{}
}

// 查找meta标签的content属性
func findMetaContent(doc *html.Node, selector map[string]string) string {
	var content string
	var traverse func(*html.Node)

	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "meta" {
			matches := true
			for key, val := range selector {
				if !hasAttribute(n, key, val) {
					matches = false
					break
				}
			}
			if matches {
				content = getAttribute(n, "content")
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if content != "" {
				return
			}
			traverse(c)
		}
	}

	traverse(doc)
	return content
}

// 查找元素的文本内容
func findElementText(doc *html.Node, tag, attr, val string) string {
	var text string
	var found bool
	var traverse func(*html.Node)

	traverse = func(n *html.Node) {
		if found {
			return
		}
		if n.Type == html.ElementNode && n.Data == tag && hasAttribute(n, attr, val) {
			text = getTextContent(n)
			found = true
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if found {
				return
			}
			traverse(c)
		}
	}

	traverse(doc)
	return text
}

// 查找元素的文本内容（通过class）
func findElementTextByClass(doc *html.Node, tag, class string) string {
	var text string
	var found bool
	var traverse func(*html.Node)

	traverse = func(n *html.Node) {
		if found {
			return
		}
		if n.Type == html.ElementNode && n.Data == tag {
			classAttr := getAttribute(n, "class")
			if strings.Contains(classAttr, class) {
				text = getTextContent(n)
				found = true
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if found {
				return
			}
			traverse(c)
		}
	}

	traverse(doc)
	return text
}

// 获取节点的文本内容
func getTextContent(n *html.Node) string {
	var buf bytes.Buffer
	var traverse func(*html.Node)

	traverse = func(node *html.Node) {
		if node.Type == html.TextNode {
			buf.WriteString(node.Data)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)
	return strings.TrimSpace(buf.String())
}

// 提取文章内容
func extractContent(doc *html.Node) (string, error) {
	var contentNode *html.Node

	// 优先级：#js_content → #img-content → .rich_media_content
	var findContent func(*html.Node)
	findContent = func(n *html.Node) {
		if contentNode != nil {
			return
		}
		if n.Type == html.ElementNode {
			// 优先从 js_content 获取
			if hasAttribute(n, "id", "js_content") {
				contentNode = n
				return
			}
			// 备选：从 img_content 获取（老版本）
			if hasAttribute(n, "id", "img-content") {
				contentNode = n
				return
			}
			// 备选：从 rich_media_content 获取
			if hasAttribute(n, "class", "rich_media_content") {
				contentNode = n
				return
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findContent(c)
		}
	}

	findContent(doc)

	if contentNode == nil {
		return "", fmt.Errorf("未找到文章内容区域")
	}

	// 清理内容节点
	cleanNode(contentNode)

	// 处理图片
	processImages(contentNode)

	// 提取文本内容
	var buf bytes.Buffer
	if err := html.Render(&buf, contentNode); err != nil {
		return "", fmt.Errorf("渲染 HTML 失败: %v", err)
	}

	return buf.String(), nil
}

// 清理节点内容
func cleanNode(n *html.Node) {
	// 移除 script 标签
	removeElements(n, "script")
	// 移除 style 标签
	removeElements(n, "style")
	// 移除注释
	removeComments(n)
	// 移除特定的微信元素
	removeSpecificElements(n)
	// 清理透明文字（反爬虫技术）
	removeTransparentText(n)
}

// 移除指定标签的所有元素
func removeElements(n *html.Node, tag string) {
	var elements []*html.Node
	var traverse func(*html.Node)

	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == tag {
			elements = append(elements, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)

	for _, e := range elements {
		removeNode(e)
	}
}

// 移除注释节点
func removeComments(n *html.Node) {
	var comments []*html.Node
	var traverse func(*html.Node)

	traverse = func(n *html.Node) {
		if n.Type == html.CommentNode {
			comments = append(comments, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)

	for _, c := range comments {
		removeNode(c)
	}
}

// 移除特定的微信元素
func removeSpecificElements(n *html.Node) {
	var elements []*html.Node
	var traverse func(*html.Node)

	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			// 移除公众号介绍相关元素
			if (n.Data == "div" && hasAttribute(n, "id", "meta_content")) ||
				(n.Data == "div" && hasAttribute(n, "id", "js_tags")) ||
				(n.Data == "div" && hasAttribute(n, "id", "js_novel_card")) ||
				(hasClass(n, "original_area_primary")) ||
				(hasClass(n, "wx_profile_card_inner")) ||
				(hasClass(n, "wx_profile_msg_inner")) {
				elements = append(elements, n)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)

	for _, e := range elements {
		removeNode(e)
	}
}

// 移除透明文字（反爬虫技术）
func removeTransparentText(n *html.Node) {
	var elements []*html.Node
	var traverse func(*html.Node)

	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			style := getAttribute(n, "style")
			if strings.Contains(style, "color: rgba(255, 255, 255, 0)") ||
				strings.Contains(style, "color: rgba(255 255 255 0)") {
				elements = append(elements, n)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)

	for _, e := range elements {
		removeNode(e)
	}
}

// 处理图片标签
func processImages(n *html.Node) {
	var traverse func(*html.Node)

	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			// 如果图片有 data-src 属性，将其赋值给 src
			dataSrc := getAttribute(n, "data-src")
			if dataSrc != "" {
				setAttribute(n, "src", dataSrc)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)
}

// 处理内容
func processContent(content string) string {
	// 去除隐藏的内容（rgba(255,255,255,0) 颜色的文字）
	re := regexp.MustCompile(`<[^>]*style\s*=\s*"[^"]*color\s*:\s*rgba\(\s*255\s*[,|\s]\s*255\s*[,|\s]\s*255\s*[,|\s]\s*0\s*\)"[^>]*>.*?</[^>]*>`)
	content = re.ReplaceAllString(content, "")

	// 移除空的段落
	emptyPara := regexp.MustCompile(`<p>\s*</p>`)
	content = emptyPara.ReplaceAllString(content, "")

	// 去除多余的换行
	content = strings.ReplaceAll(content, "\n\n", "\n")
	content = strings.TrimSpace(content)

	return content
}

// 检查节点是否有指定属性
func hasAttribute(n *html.Node, key, value string) bool {
	for _, a := range n.Attr {
		if a.Key == key && a.Val == value {
			return true
		}
	}
	return false
}

// 检查节点是否有指定class
func hasClass(n *html.Node, class string) bool {
	classAttr := getAttribute(n, "class")
	return strings.Contains(classAttr, class)
}

// 获取节点属性
func getAttribute(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

// 设置节点属性
func setAttribute(n *html.Node, key, value string) {
	found := false
	for i, a := range n.Attr {
		if a.Key == key {
			n.Attr[i].Val = value
			found = true
			break
		}
	}

	if !found {
		n.Attr = append(n.Attr, html.Attribute{Key: key, Val: value})
	}
}

// 移除节点
func removeNode(n *html.Node) {
	if n.Parent != nil {
		n.Parent.RemoveChild(n)
	}
}

// ParseHTMLFromURL 从 URL 解析 HTML 内容
func ParseHTMLFromURL(url string) ([]byte, error) {
	// 创建自定义的 HTTP 客户端
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 检查是否已经超过了最大重定向次数
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}

			// 为后续请求添加适当的头信息
			req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1")
			req.Header.Set("Referer", url)

			return nil
		},
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 模拟微信内置浏览器的请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", url)
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 URL 失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == 302 {
			return nil, fmt.Errorf("请求被重定向，请检查是否需要验证")
		}
		return nil, fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	// 处理压缩响应
	var reader io.Reader = resp.Body
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(reader)
		if err != nil {
			return nil, fmt.Errorf("解压缩 gzip 响应失败: %v", err)
		}
		defer reader.(io.ReadCloser).Close()
	case "deflate":
		reader = flate.NewReader(reader)
		defer reader.(io.ReadCloser).Close()
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("读取响应内容失败: %v", err)
	}

	return body, nil
}

// ParseFromURL 从 URL 解析微信文章
func ParseFromURL(url string) (*WechatArticle, error) {
	// 获取 HTML 内容
	htmlContent, err := ParseHTMLFromURL(url)
	if err != nil {
		return nil, err
	}

	// 解析文章
	return ParseWechatArticle(string(htmlContent), url)
}
