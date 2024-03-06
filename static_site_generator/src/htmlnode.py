from __future__ import annotations


class HTMLNode:
    def __init__(
        self,
        tag: str = "",
        value: str = "",
        children: list[HTMLNode] = [],
        props: dict[str, str] = {},
    ) -> None:
        self.tag = tag
        self.value = value
        self.children = children
        self.props = props

    def to_html(self) -> str:
        raise NotImplemented

    def props_to_html(self) -> str:
        out = " "
        for k, v in self.props.items():
            out += k + '="' + v + '" '

        return out[:-1]

    # TODO: Im assuming that this won't look nice. make it look nice
    def __repr__(self):
        out = f"HTMLNode({self.tag}, {self.value}, {self.props}, ["

        for child in self.children:
            out += "\n\t" + child.__repr__()

        out += "])"
        return out

    # FIXME: This is not recursively checking children for equality
    def __eq__(self, other: HTMLNode) -> bool:
        return (
            self.tag == other.tag
            and self.value == other.value
            and self.children == other.children
            and self.props == other.props
        )


class LeafNode(HTMLNode):
    def __init__(self, value: str, tag: str = "", props: dict[str, str] = {}) -> None:
        super().__init__(tag=tag, value=value, props=props)

    def to_html(self):
        if not self.value:
            raise ValueError("LeafNode must have a value.")

        if not self.tag:
            return self.value

        return f"<{self.tag}{self.props_to_html()}>{self.value}</{self.tag}>"

    # TODO: Im assuming that this won't look nice. make it look nice
    def __repr__(self):
        return f"LeafNode({self.tag}, {self.value}, {self.props})"


class ParentNode(HTMLNode):
    def __init__(
        self, tag: str = "", children: list[HTMLNode] = [], props: dict[str, str] = {}
    ) -> None:
        super().__init__(tag=tag, children=children, props=props)

    def to_html(self) -> str:
        if not self.tag:
            raise ValueError("ParentNode must have a tag")

        out = f"<{self.tag}{self.props_to_html()}>"
        for node in self.children:
            out += node.to_html()
        out += f"</{self.tag}>"

        return out

    # TODO: Im assuming that this won't look nice. make it look nice
    def __repr__(self):
        out = f"ParentNode({self.tag}, {self.value}, {self.props}, ["
        for child in self.children:
            out += "\n\t" + child.__repr__()

        out += "])"
        return out
