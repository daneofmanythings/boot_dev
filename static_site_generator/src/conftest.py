import pytest
from htmlnode import ParentNode
from textnode import TextNode


@pytest.fixture
def default_node():
    return TextNode(
        "test",
        "text",
        "test.gov",
    )


@pytest.fixture
def text_node():
    return TextNode("text", "text")


@pytest.fixture
def bold_node():
    return TextNode("bold", "bold")


@pytest.fixture
def italic_node():
    return TextNode("italic", "italic")


@pytest.fixture
def code_node():
    return TextNode("code", "code")


@pytest.fixture
def link_node():
    return TextNode("link", "link", "www.google.com")


@pytest.fixture
def image_node():
    return TextNode("image", "image", "www.google.com/image.png")


@pytest.fixture
def html_root():
    return ParentNode(tag="div", children=[], props={})
