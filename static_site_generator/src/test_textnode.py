import pytest
from textnode import LeafNode, TextNode, text_node_to_html_node


class TestTextNode:
    @pytest.fixture
    def default_node(self):
        return TextNode(
            "test",
            "test",
            "test.gov",
        )

    @pytest.fixture
    def text_node(self):
        return TextNode("text", "text")

    @pytest.fixture
    def bold_node(self):
        return TextNode("bold", "bold")

    @pytest.fixture
    def italic_node(self):
        return TextNode("italic", "italic")

    @pytest.fixture
    def code_node(self):
        return TextNode("code", "code")

    @pytest.fixture
    def link_node(self):
        return TextNode("link", "link", "www.google.com")

    @pytest.fixture
    def image_node(self):
        return TextNode("image", "image", "www.google.com/image.png")

    def test_text_node_eq(self, default_node):
        node_text = TextNode("not test", "test", "test.gov")
        node_text_type = TextNode("test", "not test", "test.gov")
        node_url = TextNode("test", "test", "test.com")
        node_no_url = TextNode("test", "test")
        assert (
            default_node != node_text
            and default_node != node_text_type
            and default_node != node_url
            and default_node != node_no_url
            and node_no_url.url == ""
        )

    def test_text_node_conversions_text(self, text_node):
        assert text_node_to_html_node(text_node) == LeafNode(value=text_node.text)

    def test_text_node_conversions_bold(self, bold_node):
        assert text_node_to_html_node(bold_node) == LeafNode(
            tag="b", value=bold_node.text
        )

    def test_text_node_conversions_italic(self, italic_node):
        assert text_node_to_html_node(italic_node) == LeafNode(
            tag="i", value=italic_node.text
        )

    def test_text_node_conversions_code(self, code_node):
        assert text_node_to_html_node(code_node) == LeafNode(
            tag="code", value=code_node.text
        )

    def test_text_node_conversions_link(self, link_node):
        assert text_node_to_html_node(link_node) == LeafNode(
            tag="a", value=link_node.text, props={"href": link_node.url}
        )

    def test_text_node_conversions_image(self, image_node):
        assert text_node_to_html_node(image_node) == LeafNode(
            tag="img",
            value="",
            props={"src": image_node.url, "alt": image_node.text},
        )
