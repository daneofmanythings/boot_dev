import pytest
from htmlnode import HTMLNode, LeafNode, ParentNode


@pytest.fixture
def default_htmlnode():
    return HTMLNode("a", "link", [], {"href": "google.com", "test": "value"})


@pytest.fixture
def anchor_leafnode():
    return LeafNode(
        tag="a", value="link", props={"href": "google.com", "test": "value"}
    )


@pytest.fixture
def paragraph_leafnode():
    return LeafNode(tag="p", value="paragraph")


@pytest.fixture
def plain_leafnode():
    return LeafNode(value="plain")


@pytest.fixture
def italic_leafnode():
    return LeafNode(tag="i", value="italic text")


@pytest.fixture
def bold_leafnode_with_prop():
    return LeafNode(tag="b", value="bold text", props={"test": "value"})


@pytest.fixture
def sample_parentnode(italic_leafnode, bold_leafnode_with_prop, plain_leafnode):
    return ParentNode(
        tag="p", children=[italic_leafnode, bold_leafnode_with_prop, plain_leafnode]
    )


@pytest.fixture
def empty_parentnode_with_props():
    return ParentNode(tag="p", children=[], props={"test": "value"})


@pytest.fixture
def singly_nested_parentnodes(sample_parentnode, empty_parentnode_with_props):
    return ParentNode(
        tag="p", children=[sample_parentnode, empty_parentnode_with_props]
    )


@pytest.fixture()
def doubley_nested_parentnodes(
    singly_nested_parentnodes, sample_parentnode, empty_parentnode_with_props
):
    return ParentNode(
        tag="a",
        children=[
            singly_nested_parentnodes,
            sample_parentnode,
            empty_parentnode_with_props,
        ],
    )


class TestHTMLNode:
    def test_repr(self, default_htmlnode):
        assert (
            default_htmlnode.__repr__()
            == "HTMLNode(a, link, {'href': 'google.com', 'test': 'value'}, [])"
        )

    def test_props_to_html(self, default_htmlnode):
        assert default_htmlnode.props_to_html() == ' href="google.com" test="value"'


class TestLeafNode:
    def test_repr(self, anchor_leafnode):
        assert (
            anchor_leafnode.__repr__()
            == "LeafNode(a, link, {'href': 'google.com', 'test': 'value'})"
        )

    def test_to_html(self, anchor_leafnode, paragraph_leafnode, plain_leafnode):
        assert (
            anchor_leafnode.to_html() == '<a href="google.com" test="value">link</a>'
            and paragraph_leafnode.to_html() == "<p>paragraph</p>"
            and plain_leafnode.to_html() == "plain"
        )


class TestParentNode:
    def test_to_html(
        self,
        sample_parentnode,
        empty_parentnode_with_props,
        singly_nested_parentnodes,
        doubley_nested_parentnodes,
    ):
        assert (
            sample_parentnode.to_html()
            == '<p><i>italic text</i><b test="value">bold text</b>plain</p>'
            and empty_parentnode_with_props.to_html() == '<p test="value"></p>'
            and singly_nested_parentnodes.to_html()
            == f"<{singly_nested_parentnodes.tag}>"
            + sample_parentnode.to_html()
            + empty_parentnode_with_props.to_html()
            + f"</{singly_nested_parentnodes.tag}>"
            and doubley_nested_parentnodes.to_html()
            == f"<{doubley_nested_parentnodes.tag}>"
            + singly_nested_parentnodes.to_html()
            + sample_parentnode.to_html()
            + empty_parentnode_with_props.to_html()
            + f"</{doubley_nested_parentnodes.tag}>"
        )
