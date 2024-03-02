from __future__ import annotations
from enum import StrEnum, Enum


class BlockType(StrEnum):
    """
    The starting strings of markdown blocks
    """

    PARAGRAPH = ""
    UNORD_LIST = ""
    HEADING1 = "#"
    HEADING2 = "##"
    HEADING3 = "###"
    HEADING4 = "####"
    HEADING5 = "#####"
    HEADING6 = "######"
    CODE = "```"
    QUOTE = ">"
    STAR_LIST = "*"
    DASH_LIST = "-"
    ORD_LIST = "1."

    @staticmethod
    def validate_start_end(block: str, block_type: BlockType) -> bool:
        return block.endswith(block_type.value)

    @staticmethod
    def validate_every(block: str, block_type: BlockType) -> bool:
        """
        This should never be given PARAGRAPH or UNORD_LIST as the block type.
        """
        block_lines = block.split("\n")
        if block_type != BlockType.ORD_LIST:
            for line in block_lines:
                if not line.startswith(block_type.value + " "):
                    return False
        else:
            for i in range(len(block_lines)):
                if not block_lines[i].startswith(str(i + 1) + ". "):
                    return False

        return True


class BlockTypeType(Enum):
    """
    These enums classify BlockTypes based on how they should be checked for validity.
    START_OF_BLOCK: only at the beginning of the block
    EVERY_LINE: at the beginning and after every new line of the block
    """

    START_OF_BLOCK = {
        BlockType.HEADING1,
        BlockType.HEADING2,
        BlockType.HEADING3,
        BlockType.HEADING4,
        BlockType.HEADING5,
        BlockType.HEADING6,
    }
    START_AND_END = {
        BlockType.CODE,
        BlockType.QUOTE,
    }
    EVERY_LINE = {
        BlockType.STAR_LIST,
        BlockType.DASH_LIST,
        BlockType.ORD_LIST,
    }

    @staticmethod
    def get_type(block_type: BlockType) -> BlockTypeType:
        if block_type in BlockTypeType.START_OF_BLOCK.value:
            return BlockTypeType.START_OF_BLOCK
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


def fetch_block_type_to_test(line_leader: str) -> BlockType:
    try:
        # Checking if there is a possible BlockType match.
        # NOTE: This can never return PARAGRAPH or UNORD_LIST as long as whitespace is stripped from
        # markdown_to_blocks
        block_type = BlockType(line_leader)
    except ValueError:
        # no match. return PARAGRAPH
        block_type = BlockType.PARAGRAPH

    return block_type


def block_to_block_type(block: str) -> BlockType:
    """
    Determines the block type of a block string.
    """
    line_leader = block.split(" ", 1)[0]
    block_type = fetch_block_type_to_test(line_leader)

    block_type_type = BlockTypeType.get_type(block_type)
    match block_type_type:
        case BlockTypeType.START_OF_BLOCK:
            # The validation process was handled because of the identification step in the try block
            return block_type
        case BlockTypeType.EVERY_LINE:
            if not BlockType.validate_every(block, block_type):
                return BlockType.PARAGRAPH
        case BlockTypeType.START_AND_END:
            if not BlockType.validate_start_end(block, block_type):
                return BlockType.PARAGRAPH
        case _:
            raise ValueError(f"Unhandled BlockTypeType: {block_type_type}")

    # HTML doesn't distinguish between the two representations of unordered lists.
    if block_type == BlockType.STAR_LIST or block_type == BlockType.DASH_LIST:
        return BlockType.UNORD_LIST

    return block_type
