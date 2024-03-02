from blocks import (
    BlockType,
    BlockTypeType,
    block_to_block_type,
    fetch_block_type_to_test,
    markdown_to_blocks,
)


class TestBlockTypeTypes:
    def test_get_type_start(self):
        block_type = BlockType.HEADING6

        assert BlockTypeType.get_type(block_type) == BlockTypeType.START_OF_BLOCK

    def test_get_type_every(self):
        block_type = BlockType.ORD_LIST

        assert BlockTypeType.get_type(block_type) == BlockTypeType.EVERY_LINE

    def test_get_type_end(self):
        block_type = BlockType.CODE

        assert BlockTypeType.get_type(block_type) == BlockTypeType.START_AND_END


class TestFetchBlockTypeToTest:
    def test_h1(self):
        assert fetch_block_type_to_test("#") == BlockType.HEADING1

    def test_out(self):
        assert fetch_block_type_to_test("a") == BlockType.PARAGRAPH

    def test_code(self):
        assert fetch_block_type_to_test("```") == BlockType.CODE


class TestBlockToBlockType:
    def test_default(self):
        text = "This is just some normal text.\nThere is nothing to see here."
        block_type = block_to_block_type(text)

        assert block_type == BlockType.PARAGRAPH

    def test_heading1(self):
        text = "# This is a HUGE heading."
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
        text = "```This is a coooode block\nAnd it is very codelike```"
        block_type = block_to_block_type(text)

        assert block_type == BlockType.CODE


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
