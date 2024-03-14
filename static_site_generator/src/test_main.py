import pytest
from blocks import markdown_to_html_node


class TestMarkdownToHTML:
    def test_only_heading(self):
        md = "# This is a heading"
        expected = "<div><h1>This is a heading</h1></div>"
        result = markdown_to_html_node(md).to_html()

        assert expected == result

    def test_only_image(self):
        md = "![image](www.image.com)"
        expected = "<div><"
