from __future__ import annotations
from enum import StrEnum, Enum


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
    # TODO:: fix star and dash to mix in lists with no issue
    STAR_LIST = "*"
    DASH_LIST = "-"
    UNORD_LIST = ""

    PARAGRAPH = ""

    @staticmethod
    def validate_block_end(block: str, block_type: BlockType) -> bool:
        # NOTE: We got here by the block leading with "```". Now we check the splits length
        # to see if there are the correct number of surrounds
        # TODO: Fix this for when someone gets triggerhappy with backticks
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
        else:
            for line in block_lines:
                if not line.startswith(block_type.value):
                    return False

        return True


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
    }
    START_AND_END = {
        BlockType.CODE,
    }
    EVERY_LINE = {
        BlockType.QUOTE,
        BlockType.STAR_LIST,
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


def peek_block_type(block: str) -> BlockType:
    # NOTE: this is clunkier than i'd like.
    for block_type_enum in BlockType.__members__.values():
        if not block_type_enum.value:
            continue
        if block.startswith(block_type_enum.value):
            return block_type_enum

    return BlockType.PARAGRAPH


def block_to_block_type(block: str) -> BlockType:
    """
    Determines the block type of a block string.
    """
    # NOTE: This should never peek a "" because we have stripped all whitespace
    block_type = peek_block_type(block)

    block_type_type = BlockTypeType.get_type(block_type)
    match block_type_type:
        case BlockTypeType.START:
            # The validation process was handled because of the identification
            # step in the try block
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
    if block_type == BlockType.STAR_LIST or block_type == BlockType.DASH_LIST:
        return BlockType.UNORD_LIST

    return block_type
