from textnode import TextNode


def main():
    text = "This is some text"
    text_type = "bold"
    url = "google.com"

    tn = TextNode(text, text_type, url)

    print(tn)


if __name__ == "__main__":
    main()
