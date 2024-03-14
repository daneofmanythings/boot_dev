from __future__ import annotations
from enum import StrEnum, Enum
import re
from htmlnode import HTMLNode, LeafNode, ParentNode
from textnode import text_node_to_html_node
from text_to_textnode import text_to_textnode


class BlockType(StrEnum):
    """
    The starting strings of markdown blocks. The declaration order matters for the parsing.
    """

    HEADING6 = "###### "
    HEADING5 = "##### "
    HEADING4 = "#### "
    HEADING3 = "### "
    HEADING2 = "## "
    HEADING1 = "# "
    CODE = "```"
    QUOTE = ">"
    ORD_LIST = "1."
    UNORD_LIST = "*"
    DASH_LIST = "-"
    PARAGRAPH = ""

    @staticmethod
    def validate_block_end(block: str, block_type: BlockType) -> bool:
        # NOTE: We got here by the block leading with "```". Now we check the splits length
        # to see if there are the correct number of surrounds
        # TODO: Fix this for when someone gets triggerhappy with backticks for codeblocks
        return len(block.split(block_type.value)) % 2 == 1

    @staticmethod
    def validate_every_line(block: str, block_type: BlockType) -> bool:
        """
        This should never be given PARAGRAPH or UNORD_LIST as the block type.
        """
        block_lines = block.split("\n")
        if block_type == BlockType.ORD_LIST:
            for i in range(len(block_lines)):
                if not block_lines[i].startswith(str(i + 1) + "."):
                    return False
        # NOTE: verifies the unordered list regardless of the delimiter
        elif block_type == BlockType.UNORD_LIST or block_type == BlockType.DASH_LIST:
            for line in block_lines:
                if not is_unordered(line):
                    print(f"block type val failed for block: {block} on line: {line}")
                    return False
        else:
            for line in block_lines:
                if not line.startswith(block_type.value):
                    return False

        return True


def is_unordered(line: str) -> bool:
    unord = line.startswith(BlockType.UNORD_LIST.value)
    dash = line.startswith(BlockType.DASH_LIST.value)
    return unord or dash


class BlockTypeType(Enum):
    """
    These enums classify BlockTypes based on how they should be checked for validity.
    START: only at the beginning of the block
    START_AND_END: at the start and again at the end of the block
    EVERY_LINE: at the beginning and after every new line of the block
    """

    START = {
        BlockType.HEADING1,
        BlockType.HEADING2,
        BlockType.HEADING3,
        BlockType.HEADING4,
        BlockType.HEADING5,
        BlockType.HEADING6,
        BlockType.PARAGRAPH,
    }
    START_AND_END = {
        BlockType.CODE,
    }
    EVERY_LINE = {
        BlockType.QUOTE,
        BlockType.UNORD_LIST,
        BlockType.DASH_LIST,
        BlockType.ORD_LIST,
    }

    @staticmethod
    def get_type(block_type: BlockType) -> BlockTypeType:
        if block_type in BlockTypeType.START.value:
            return BlockTypeType.START
        if block_type in BlockTypeType.START_AND_END.value:
            return BlockTypeType.START_AND_END
        return BlockTypeType.EVERY_LINE


def markdown_to_blocks(markdown: str) -> list[str]:
    """
    Splits the markdown string into a list of stripped strings.
    """
    block_delimeter = "\n\n"

    blocks = markdown.split(block_delimeter)
    for i in range(len(blocks)):
        blocks[i] = blocks[i].strip()

    blocks = list(filter(lambda b: b, blocks))

    return blocks


# NOTE: this is clunkier than i'd like.
def peek_block_type(block: str) -> BlockType:
    # loops through all blocktypes to find a match
    for block_type_enum in BlockType.__members__.values():
        if not block_type_enum.value:  # skips the "" case.
            continue
        if block.startswith(block_type_enum.value):  # returns a match
            return block_type_enum

    return BlockType.PARAGRAPH  # returns if there is no match


def block_to_block_type(block: str) -> BlockType:
    """
    Determines the block type of a block string.
    """
    # NOTE: This should never peek a "" because we have stripped all whitespace
    block_type = peek_block_type(block)

    block_type_type = BlockTypeType.get_type(block_type)
    match block_type_type:
        case BlockTypeType.START:
            return block_type

        case BlockTypeType.START_AND_END:
            if not BlockType.validate_block_end(block, block_type):
                return BlockType.PARAGRAPH

        case BlockTypeType.EVERY_LINE:
            if not BlockType.validate_every_line(block, block_type):
                return BlockType.PARAGRAPH

        case _:
            raise ValueError(f"Unhandled BlockTypeType: {block_type_type}")

    # HTML doesn't distinguish between the two representations of unordered lists.
    if block_type == BlockType.DASH_LIST:
        return BlockType.UNORD_LIST

    return block_type


def heading_to_htmlnode(block: str, block_type: BlockType) -> HTMLNode:
    header_len = len(block_type.value)
    tag = f"h{header_len - 1}"
    textnodes = text_to_textnode(block[header_len:])
    children = list()
    for node in textnodes:
        children.append(text_node_to_html_node(node))

    return ParentNode(tag=tag, children=children, props={})


def code_to_html(block: str) -> HTMLNode:
    # All code nodes are packed into a <pre>
    pre_node = ParentNode(tag="pre", children=[], props={})
    tag = "code"
    # Isolating code blocks and stripping out any that might not be
    code_blocks = block.split("```")
    code_blocks = filter(lambda x: x and not x.isspace(), code_blocks)
    # For each code block, we ...
    for code_node in code_blocks:
        # ... make its textnodes ...
        textnodes = text_to_textnode(code_node)
        children = list()
        # ... convert them and add them as children ...
        for textnode in textnodes:
            children.append(text_node_to_html_node(textnode))
        # ... make and append the code node to the pre node.
        pre_node.children.append(ParentNode(tag=tag, children=children, props={}))

    return pre_node


def quote_to_html(block: str) -> HTMLNode:
    tag = f"blockquote"
    text = "\n".join(block.split("\n>"))
    textnodes = text_to_textnode(text[1:])
    children = list()
    for node in textnodes:
        children.append(text_node_to_html_node(node))

    return ParentNode(tag=tag, children=children, props={})


def ord_list_to_html(block: str) -> HTMLNode:
    # Setup for the parent node
    tag = "ol"
    children = list()

    # Splitting up the list items to be thier own nodes
    list_items = re.split(r"\n\d+\.", block[2:])  # Cropping out first delimiter

    # Constructing the list item nodes
    for item in list_items:
        children.append(list_item_to_html(item))

    return ParentNode(tag=tag, children=children, props={})


def unord_list_to_html(block: str) -> HTMLNode:
    # Setup for the parent node
    tag = "ul"
    children = list()

    # Normalizing the list delimiters
    block = block[1:].replace("\n*", "\n-")  # Cropping out the first delimiter

    # Splitting up the list items to be thier own nodes
    list_items = block.split("\n-")

    # Constructing the list item nodes
    for item in list_items:
        children.append(list_item_to_html(item))

    return ParentNode(tag=tag, children=children, props={})


def list_item_to_html(list_item: str) -> HTMLNode:
    tag = "li"
    children = list()

    textnodes = text_to_textnode(list_item)

    for node in textnodes:
        children.append(text_node_to_html_node(node))

    return ParentNode(tag=tag, children=children, props={})


def paragraph_to_html(block: str) -> HTMLNode:
    tag = "p"
    children = list()

    textnodes = text_to_textnode(block)

    for node in textnodes:
        children.append(text_node_to_html_node(node))

    return ParentNode(tag=tag, children=children, props={})


def markdown_to_html_node(markdown: str) -> HTMLNode:
    blocks = markdown_to_blocks(markdown)
    root = ParentNode(tag="div", children=[], props={})
    for block in blocks:
        block_type = block_to_block_type(block)
        print(block_type.name)
        match block_type:
            case (
                BlockType.HEADING6
                | BlockType.HEADING5
                | BlockType.HEADING4
                | BlockType.HEADING3
                | BlockType.HEADING2
                | BlockType.HEADING1
            ):
                block_html = heading_to_htmlnode(block, block_type)
            case BlockType.CODE:
                block_html = code_to_html(block)
            case BlockType.QUOTE:
                block_html = quote_to_html(block)
            case BlockType.ORD_LIST:
                block_html = ord_list_to_html(block)
            case BlockType.UNORD_LIST:
                block_html = unord_list_to_html(block)
            case BlockType.PARAGRAPH:
                block_html = paragraph_to_html(block)
            case _:
                raise ValueError(f"Unhandled BlockType: {block_type}")

        root.children.append(block_html)

    return root
