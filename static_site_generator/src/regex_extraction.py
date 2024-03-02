import re


def extract_markdown_images(text: str) -> list[tuple[str, str]]:
    result = list()
    matches = re.findall(r"!\[(.*?)\]\((.*?)\)", text)
    result.extend(matches)

    return result


def extract_markdown_links(text: str) -> list[tuple[str, str]]:
    result = list()
    matches = re.findall(r"\[(.*?)\]\((.*?)\)", text)
    result.extend(matches)

    return result
