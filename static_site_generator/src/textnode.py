from __future__ import annotations

from htmlnode import HTMLNode, LeafNode


class TextNode:
    def __init__(self, text: str, text_type: str, url: str = "") -> None:
        self.text = text
        self.text_type = text_type
        self.url = url

    def __eq__(self, other: TextNode) -> bool:
        return (
            self.text == other.text
            and self.text_type == other.text_type
            and self.url == other.url
        )

    def __repr__(self) -> str:
        return f"TextNode({self.text}, {self.text_type}, {self.url})"


def text_node_to_html_node(text_node: TextNode) -> HTMLNode:
    match text_node.text_type:
        case "text":
            return text_type_text(text_node)
        case "bold":
            return text_type_bold(text_node)
        case "italic":
            return text_type_italic(text_node)
        case "code":
            return text_type_code(text_node)
        case "link":
            return text_type_link(text_node)
        case "image":
            return text_type_image(text_node)
        case _:
            raise ValueError(
                f"TextNode must have a type of 'text', 'bold', 'italic', 'code', 'link', 'image'. got={text_node.text_type}"
            )


def text_type_text(text_node: TextNode) -> HTMLNode:
    return LeafNode(value=text_node.text)


def text_type_bold(text_node: TextNode) -> HTMLNode:
    return LeafNode(tag="b", value=text_node.text)


def text_type_italic(text_node: TextNode) -> HTMLNode:
    return LeafNode(tag="i", value=text_node.text)


def text_type_code(text_node: TextNode) -> HTMLNode:
    return LeafNode(tag="code", value=text_node.text)


def text_type_link(text_node: TextNode) -> HTMLNode:
    return LeafNode(tag="a", value=text_node.text, props={"href": text_node.url})


def text_type_image(text_node: TextNode) -> HTMLNode:
    return LeafNode(
        tag="img", value="", props={"src": text_node.url, "alt": text_node.text}
    )
