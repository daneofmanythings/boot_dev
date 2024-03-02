from regex_extraction import extract_markdown_images, extract_markdown_links


class TestExtractMarkdownImages:
    def test_sample(self):
        text = "This is text with an ![image](https://i.imgur.com/zjjcJKZ.png) and ![another](https://i.imgur.com/dfsdkjfd.png)"
        matches = extract_markdown_images(text)
        assert matches == [
            ("image", "https://i.imgur.com/zjjcJKZ.png"),
            ("another", "https://i.imgur.com/dfsdkjfd.png"),
        ]

    def test_empty(self):
        text = ""
        matches = extract_markdown_images(text)
        assert matches == []

    def test_no_matches(self):
        text = "this text has ![no] (matches)"
        matches = extract_markdown_images(text)
        assert matches == []

    def test_only_matches(self):
        text = "![image](https://i.imgur.com/zjjcJKZ.png)![another](https://i.imgur.com/dfsdkjfd.png)"
        matches = extract_markdown_images(text)
        assert matches == [
            ("image", "https://i.imgur.com/zjjcJKZ.png"),
            ("another", "https://i.imgur.com/dfsdkjfd.png"),
        ]


class TestExtractMarkdownLinks:
    def test_sample(self):
        text = "This is text with a [link](https://www.example.com) and [another](https://www.example.com/another)"
        matches = extract_markdown_links(text)
        assert matches == [
            ("link", "https://www.example.com"),
            ("another", "https://www.example.com/another"),
        ]

    def test_empty(self):
        text = ""
        matches = extract_markdown_links(text)
        assert matches == []

    def test_no_matches(self):
        text = "this text has [no] (matches)"
        matches = extract_markdown_links(text)
        assert matches == []

    def test_only_matches(self):
        text = (
            "[link](https://www.example.com)[another](https://www.example.com/another)"
        )
        matches = extract_markdown_links(text)
        assert matches == [
            ("link", "https://www.example.com"),
            ("another", "https://www.example.com/another"),
        ]
