from textnode import TextNode, TextTypeNode
from text_to_textnode import (
    split_nodes_delimiter,
    split_nodes_image,
    split_nodes_link,
    text_to_textnode,
)


class TestSplitNodesDelimiter:
    def test_sample(self):
        test_text_node = TextNode("test `code block` test", "text")
        assert split_nodes_delimiter([test_text_node], "`", TextTypeNode.CODE) == [
            TextNode("test ", "text"),
            TextNode("code block", "code"),
            TextNode(" test", "text"),
        ]

    def test_empty(self):
        assert split_nodes_delimiter([], "*", TextTypeNode.ITALIC) == []

    def test_multi(self, code_node, link_node):
        test_text_node = TextNode("test **bold** test", "text")
        old_nodes = [code_node, test_text_node, link_node]

        assert split_nodes_delimiter(old_nodes, "**", TextTypeNode.BOLD) == [
            code_node,
            TextNode("test ", "text"),
            TextNode("bold", "bold"),
            TextNode(" test", "text"),
            link_node,
        ]

    def test_double(self, image_node):
        test_text_node = TextNode("this is **bold** and this is **bold**", "text")
        old_nodes = [image_node, test_text_node]

        assert split_nodes_delimiter(old_nodes, "**", TextTypeNode.BOLD) == [
            image_node,
            TextNode("this is ", "text"),
            TextNode("bold", "bold"),
            TextNode(" and this is ", "text"),
            TextNode("bold", "bold"),
        ]


class TestSplitNodesImage:
    def test_simple(self):
        node = TextNode("an image: ![image](https://image.com)", "text")
        new_nodes = split_nodes_image([node])

        assert new_nodes == [
            TextNode("an image: ", "text"),
            TextNode("image", "image", "https://image.com"),
        ]

    def test_sample(self):
        node = TextNode(
            "This is text with an ![image](https://i.imgur.com/zjjcJKZ.png) and another ![second image](https://i.imgur.com/3elNhQu.png)",
            "text",
        )
        new_nodes = split_nodes_image([node])

        assert new_nodes == [
            TextNode("This is text with an ", TextTypeNode.TEXT),
            TextNode("image", TextTypeNode.IMAGE, "https://i.imgur.com/zjjcJKZ.png"),
            TextNode(" and another ", TextTypeNode.TEXT),
            TextNode(
                "second image", TextTypeNode.IMAGE, "https://i.imgur.com/3elNhQu.png"
            ),
        ]

    def test_no_image(self):
        node = TextNode("there is ![no] (image)", "text")
        new_nodes = split_nodes_image([node])

        assert new_nodes == [node]

    def test_only_image(self):
        node = TextNode("![image](https://i.imgur.com/3elNhQu.png)", "text")
        new_nodes = split_nodes_image([node])

        assert new_nodes == [
            TextNode("image", "image", "https://i.imgur.com/3elNhQu.png")
        ]

    def test_split_image_preserve_bold(self):
        node = TextNode("this text is bold", "bold")
        new_nodes = split_nodes_image([node])

        assert new_nodes == [node]


class TestSplitNodesLink:
    def test_simple(self):
        node = TextNode("a link: [link](https://image.com)", "text")
        new_nodes = split_nodes_link([node])

        assert new_nodes == [
            TextNode("a link: ", "text"),
            TextNode("link", "link", "https://image.com"),
        ]

    def test_sample(self):
        node = TextNode(
            "This is text with an [link](https://i.imgur.com/zjjcJKZ.png) and another [second link](https://i.imgur.com/3elNhQu.png)",
            "text",
        )
        new_nodes = split_nodes_link([node])

        assert new_nodes == [
            TextNode("This is text with an ", TextTypeNode.TEXT),
            TextNode("link", TextTypeNode.LINK, "https://i.imgur.com/zjjcJKZ.png"),
            TextNode(" and another ", TextTypeNode.TEXT),
            TextNode(
                "second link", TextTypeNode.LINK, "https://i.imgur.com/3elNhQu.png"
            ),
        ]

    def test_no_link(self):
        node = TextNode("there is [no] (link)", "text")
        new_nodes = split_nodes_link([node])

        assert new_nodes == [node]

    def test_only_link(self):
        node = TextNode("[link](https://i.imgur.com/3elNhQu.png)", "text")
        new_nodes = split_nodes_link([node])

        assert new_nodes == [
            TextNode("link", "link", "https://i.imgur.com/3elNhQu.png")
        ]


class TestTextToTextNode:
    def test_sample(self):
        text = "This is **text** with an *italic* word and a `code block` and an ![image](https://i.imgur.com/zjjcJKZ.png) and a [link](https://boot.dev)"

        textnodes = text_to_textnode(text)

        assert textnodes == [
            TextNode("This is ", "text"),
            TextNode("text", "bold"),
            TextNode(" with an ", "text"),
            TextNode("italic", "italic"),
            TextNode(" word and a ", "text"),
            TextNode("code block", "code"),
            TextNode(" and an ", "text"),
            TextNode("image", "image", "https://i.imgur.com/zjjcJKZ.png"),
            TextNode(" and a ", "text"),
            TextNode("link", "link", "https://boot.dev"),
        ]

    # TODO: Figure out if this test shoult even pass. is it valid md?

    # def test_no_plain(self):
    #     text = "*This is a ***string with no **`plain text. `![it is for ](testing.com)"
    #     textnodes = text_to_textnode(text)
    #
    #     assert textnodes == [
    #         TextNode("This is a ", "italic"),
    #         TextNode("string with no ", "bold"),
    #         TextNode("plain text. ", "code"),
    #         TextNode("it is for ", "image", "testing.com"),
    #     ]
