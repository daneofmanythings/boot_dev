from blocks import markdown_to_html_node
from htmlnode import HTMLNode
import os
import shutil


def main():
    static_path = "./static/"
    public_path = "./public/"
    content_path = "./content/"
    template_path = "./template.html"

    copy_static_files(static_path, public_path)
    # Generate index page
    generate_page(content_path + "index.md", template_path, public_path)


def copy_static_files(read_path, write_path):
    # remove the public directory
    if os.path.exists(write_path):
        shutil.rmtree(write_path)

    # create the public directory
    os.mkdir(write_path)

    # enter the copy recursion
    copy_static_files_r(read_path, write_path)


def copy_static_files_r(read_path: str, write_path: str):
    # get the read path's contents
    contents = os.listdir(read_path)

    # for each content ...
    for content in contents:
        # convenience variables
        cur_read = os.path.join(read_path, content)
        cur_write = os.path.join(write_path, content)

        # ... if we found a dir ...
        if os.path.isdir(cur_read):
            # ... make the dir in the write path and recurse on it ...
            os.mkdir(cur_write)
            copy_static_files_r(cur_read, cur_write)
        else:
            # ... or copy the file we found
            shutil.copy(cur_read, write_path)


def generate_page(from_path: str, template_path: str, dest_path: str):
    print(f"Generating page from {from_path} to {dest_path} using {template_path}...")
    markdown = ""
    template = ""
    title_tag = "{{ Title }}"
    content_tag = "{{ Content }}"

    with open(from_path, "r") as f:
        markdown = f.read()

    with open(template_path, "r") as f:
        template = f.read()

    html_nodes = markdown_to_html_node(markdown)

    title, html_nodes = extract_title(html_nodes)
    content = html_nodes.to_html()

    page = ""

    page = template.replace(title_tag, title)
    page = template.replace(content_tag, content)

    os.makedirs(dest_path, exist_ok=True)

    with open(dest_path + "index.html", "w") as f:
        f.write(page)


def extract_title(html: HTMLNode) -> tuple[str, HTMLNode]:
    title = "This is a placeholder"

    return title, html


if __name__ == "__main__":
    main()
