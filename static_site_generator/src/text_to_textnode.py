from textnode import TextNode, TextTypeNode
from regex_extraction import extract_markdown_images, extract_markdown_links


# NOTE: These implementations does not currently support nested delimiters
def split_nodes_delimiter(
    old_nodes: list[TextNode], delimiter: str, text_type: TextTypeNode
) -> list[TextNode]:

    new_nodes = list()

    for node in old_nodes:
        if not isinstance(node, TextNode):
            new_nodes.append(node)
            continue

        split_on_delim = node.text.split(delimiter)

        if len(split_on_delim) % 2 != 1:
            raise ValueError("Invalid markdown syntax: " + node.text)

        for i, new_node_text in enumerate(split_on_delim):
            if not new_node_text:
                continue

            if i % 2 == 0:
                new_nodes.append(TextNode(new_node_text, node.text_type, url=node.url))
            else:
                new_nodes.append(TextNode(new_node_text, text_type))

    return new_nodes


def split_nodes_image(old_nodes: list[TextNode]) -> list[TextNode]:
    new_nodes = list()

    for node in old_nodes:
        if not isinstance(node, TextNode):
            new_nodes.append(node)
            continue

        image_matches = extract_markdown_images(node.text)

        text_to_parse = node.text
        for image_info in image_matches:
            delimiter = "![" + image_info[0] + "](" + image_info[1] + ")"
            try:
                plain_text, text_to_parse = text_to_parse.split(delimiter, 1)
            except ValueError as e:
                raise ValueError("There was an error parsing image data: " + str(e))

            if plain_text:
                new_nodes.append(TextNode(plain_text, node.text_type))
            new_nodes.append(TextNode(image_info[0], "image", image_info[1]))

        if text_to_parse:
            new_nodes.append(TextNode(text_to_parse, node.text_type))

    return new_nodes


def split_nodes_link(old_nodes: list[TextNode]) -> list[TextNode]:
    new_nodes = list()

    for node in old_nodes:
        if not isinstance(node, TextNode):
            new_nodes.append(node)
            continue

        image_matches = extract_markdown_links(node.text)

        text_to_parse = node.text
        for link_info in image_matches:
            delimiter = "[" + link_info[0] + "](" + link_info[1] + ")"
            try:
                plain_text, text_to_parse = text_to_parse.split(delimiter, 1)
            except ValueError as e:
                raise ValueError("There was an error parsing link data: " + str(e))

            if plain_text:
                new_nodes.append(TextNode(plain_text, node.text_type))
            new_nodes.append(TextNode(link_info[0], "link", link_info[1]))

        if text_to_parse:
            new_nodes.append(TextNode(text_to_parse, node.text_type, node.url))

    return new_nodes


def text_to_textnode(text: str) -> list[TextNode]:
    base_node = TextNode(text, "text")
    textnodes = [base_node]

    textnodes = split_nodes_delimiter(textnodes, "**", TextTypeNode.BOLD)
    textnodes = split_nodes_delimiter(textnodes, "*", TextTypeNode.ITALIC)
    textnodes = split_nodes_delimiter(textnodes, "`", TextTypeNode.CODE)
    textnodes = split_nodes_image(textnodes)
    textnodes = split_nodes_link(textnodes)

    return textnodes
