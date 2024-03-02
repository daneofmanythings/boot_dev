from textnode import (
    LeafNode,
    TextNode,
    text_node_to_html_node,
)


class TestTextNode:

    def test_text_node_eq(self, default_node):
        node_text = TextNode("not test", "text", "test.gov")
        node_text_type = TextNode("test", "bold", "test.gov")
        node_url = TextNode("test", "text", "test.com")
        node_no_url = TextNode("test", "text")
        assert (
            default_node != node_text
            and default_node != node_text_type
            and default_node != node_url
            and default_node != node_no_url
            and node_no_url.url == ""
        )


class TestTextNodeConversions:

    def test_text(self, text_node):
        assert text_node_to_html_node(text_node) == LeafNode(value=text_node.text)

    def test_bold(self, bold_node):
        assert text_node_to_html_node(bold_node) == LeafNode(
            tag="b", value=bold_node.text
        )

    def test_italic(self, italic_node):
        assert text_node_to_html_node(italic_node) == LeafNode(
            tag="i", value=italic_node.text
        )

    def test_code(self, code_node):
        assert text_node_to_html_node(code_node) == LeafNode(
            tag="code", value=code_node.text
        )

    def test_link(self, link_node):
        assert text_node_to_html_node(link_node) == LeafNode(
            tag="a", value=link_node.text, props={"href": link_node.url}
        )

    def test_image(self, image_node):
        assert text_node_to_html_node(image_node) == LeafNode(
            tag="img",
            value="",
            props={"src": image_node.url, "alt": image_node.text},
        )
