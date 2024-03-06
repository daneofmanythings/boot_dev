from blocks import (
    BlockType,
    BlockTypeType,
    block_to_block_type,
    peek_block_type,
    markdown_to_blocks,
    markdown_to_html_node,
)
from htmlnode import ParentNode, LeafNode


class TestBlockTypes:
    def test_validate_block_end(self):
        text = "```This is a code block.```"
        block_type = BlockType.CODE

        assert BlockType.validate_block_end(text, block_type)

    def test_validate_block_end_failure(self):
        text = "```This is a code block.``"
        block_type = BlockType.CODE

        assert not BlockType.validate_block_end(text, block_type)

    def test_validate_every_line_ordered(self):
        text = "1. This\n2. is an\n3. ordered list."
        block_type = BlockType.ORD_LIST

        assert BlockType.validate_every_line(text, block_type)


class TestBlockTypeTypes:
    def test_get_type_start(self):
        block_type = BlockType.HEADING6

        assert BlockTypeType.get_type(block_type) == BlockTypeType.START

    def test_get_type_every(self):
        block_type = BlockType.ORD_LIST

        assert BlockTypeType.get_type(block_type) == BlockTypeType.EVERY_LINE

    def test_get_type_end(self):
        block_type = BlockType.CODE

        assert BlockTypeType.get_type(block_type) == BlockTypeType.START_AND_END


class TestPeekBlockType:
    def test_h1(self):
        assert peek_block_type("# ") == BlockType.HEADING1

    def test_out(self):
        assert peek_block_type("a") == BlockType.PARAGRAPH

    def test_code(self):
        assert peek_block_type("```") == BlockType.CODE


class TestBlockToBlockType:
    def test_default(self):
        text = "This is just some normal text.\nThere is nothing to see here."
        block_type = block_to_block_type(text)

        assert block_type == BlockType.PARAGRAPH

    def test_heading1(self):
        text = "# This is a heading."
        block_type = block_to_block_type(text)

        assert block_type == BlockType.HEADING1

    def test_heading5(self):
        text = "##### This is a HUGE heading."
        block_type = block_to_block_type(text)

        assert block_type == BlockType.HEADING5

    def test_malformed_heading2(self):
        text = "##This is a malformed heading"
        block_type = block_to_block_type(text)

        assert block_type == BlockType.PARAGRAPH

    def test_star_list(self):
        text = "* This is\n* a star\n* list."
        block_type = block_to_block_type(text)

        assert block_type == BlockType.UNORD_LIST

    def test_dash_list(self):
        text = "- This is\n- a star\n- list."
        block_type = block_to_block_type(text)

        assert block_type == BlockType.UNORD_LIST

    def test_malformed_unordeded_list(self):
        text = "* This is\n* a star\n- list."
        block_type = block_to_block_type(text)

        assert block_type == BlockType.PARAGRAPH

    def test_ordered_list(self):
        text = "1. This is\n2. an ordered\n3. list."
        block_type = block_to_block_type(text)

        assert block_type == BlockType.ORD_LIST

    def test_malformed_ordered_list(self):
        text = "1. This is\n3. an ordered\n3. list."
        block_type = block_to_block_type(text)

        assert block_type == BlockType.PARAGRAPH

    def test_code(self):
        text = "```\nThis is code\n``` ``` And it has multiple blocks ```"
        block_type = block_to_block_type(text)

        assert block_type == BlockType.CODE

    def test_quote(self):
        text = "> This is a quote block.\n> Don't quote me on that"
        block_type = block_to_block_type(text)

        assert block_type == BlockType.QUOTE


class TestMarkdownToBlocks:
    def test_simple(self):
        text = "This is\nsimple"
        blocks = markdown_to_blocks(text)

        assert blocks == [text]

    def test_sample(self):
        text = "This is a **bolded** paragraph\n\nThis is another paragraph with *italic* text and `code` here\nThis is the same paragraph on a new line\n\n* This is a list\n* with items"
        blocks = markdown_to_blocks(text)

        assert blocks == [
            "This is a **bolded** paragraph",
            "This is another paragraph with *italic* text and `code` here\nThis is the same paragraph on a new line",
            "* This is a list\n* with items",
        ]

    def test_whitespace(self):
        text = "\n\n\n\n   This\n\n   \nhas\na lot\n\nof\n\n  \n\nwhitespace\n  \n\n\n"
        blocks = markdown_to_blocks(text)

        assert blocks == ["This", "has\na lot", "of", "whitespace"]


class TestBlockToHTML:
    def test_header1_simple(self, html_root):
        text = "# This is a header 1."
        html_node = markdown_to_html_node(text)
        node = ParentNode(tag="h1", children=[LeafNode(value="This is a header 1.")])
        html_root.children.append(node)

        assert html_node == html_root

    def test_quoteblock_simple(self, html_root):
        text = "> This is a blockquote.\n> It has more than one line."
        html_node = markdown_to_html_node(text)
        node = ParentNode(
            tag="blockquote",
            children=[
                LeafNode(value=" This is a blockquote.\n It has more than one line.")
            ],
        )
        html_root.children.append(node)

        assert html_node == html_root

    def test_unordered_list_simple(self, html_root):
        text = "* This is \n- an unordered\n*list."
        html_node = markdown_to_html_node(text)
        node = ParentNode(
            tag="ul",
            children=[
                ParentNode(tag="li", children=[LeafNode(value=" This is ")], props={}),
                ParentNode(
                    tag="li", children=[LeafNode(value=" an unordered")], props={}
                ),
                ParentNode(tag="li", children=[LeafNode(value="list.")], props={}),
            ],
            props={},
        )
        html_root.children.append(node)

        assert html_node == html_root

    def test_ordered_list_simple(self, html_root):
        text = "1. This is \n2. an ordered\n3.list."
        html_node = markdown_to_html_node(text)
        node = ParentNode(
            tag="ol",
            children=[
                ParentNode(tag="li", children=[LeafNode(value=" This is ")], props={}),
                ParentNode(
                    tag="li", children=[LeafNode(value=" an ordered")], props={}
                ),
                ParentNode(tag="li", children=[LeafNode(value="list.")], props={}),
            ],
            props={},
        )
        html_root.children.append(node)

        assert html_node == html_root

    def test_code_simple(self, html_root):
        text = "```This is a code block.```\n```And another one```"
        html_node = markdown_to_html_node(text)
        node = ParentNode(
            tag="pre",
            children=[
                ParentNode(
                    tag="code",
                    children=[LeafNode(value="This is a code block.")],
                    props={},
                ),
                ParentNode(
                    tag="code", children=[LeafNode(value="And another one")], props={}
                ),
            ],
            props={},
        )

        html_root.children.append(node)

        assert html_node == html_root
