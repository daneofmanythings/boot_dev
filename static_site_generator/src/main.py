from blocks import markdown_to_html_node
from htmlnode import HTMLNode
import os
import pathlib
import shutil


def main():
    static_path = "./static/"
    public_path = "./public/"
    content_path = "./content/"
    template_path = "./template.html"

    copy_static_files(static_path, public_path)
    # Generate index page
    generate_pages(content_path, template_path, public_path)


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


def generate_pages(content_path_dir: str, template_path: str, dest_path_dir: str):
    contents = os.listdir(content_path_dir)

    for content in contents:
        cur_read = os.path.join(content_path_dir, content)

        if os.path.isdir(cur_read):
            cur_write = os.path.join(dest_path_dir, content)
            os.mkdir(cur_write)
            generate_pages(cur_read, template_path, cur_write)
        else:
            file_name = pathlib.Path(content).with_suffix(".html")
            generate_page(
                cur_read, template_path, os.path.join(dest_path_dir, file_name)
            )


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
    page = page.replace(content_tag, content)

    dirs_to_make = pathlib.Path(dest_path).parents
    os.makedirs(str(dirs_to_make), exist_ok=True)

    with open(dest_path, "w") as f:
        f.write(page)


def extract_title(html: HTMLNode) -> tuple[str, HTMLNode]:
    header = html.children[0]
    if header.tag != "h1":
        raise ValueError(f"Header 1 not found, got tag={header.tag} instead.")
    title = header.children[0].value

    return title, html


if __name__ == "__main__":
    main()
