import string


def char_count(file_str: str) -> dict[str:int]:
    count = dict()
    for c in file_str:
        c_ = c.lower()
        if c_ not in string.ascii_lowercase:
            continue
        if c_ not in count:
            count[c_] = 0
        count[c_] += 1
    return count


def num_words(file_str: str) -> int:
    words = file_str.split()
    return len(words)


def dict_to_tuple_list(dictionary: dict) -> list((str, int)):
    l = list()
    for k, v in dictionary.items():
        l.append((k, v))
    return l


def print_report(file_str: str):
    print('--- Begin report of books/frankenstein.txt ---')
    print(f"{num_words(file_str)} words found in the document\n")
    count = char_count(file_str)
    count_list = dict_to_tuple_list(count)
    count_list.sort(key=lambda x: x[1])
    for x in count_list:
        print(f"The '{x[0]} character was found {x[1]} times'")

    print('--- End report ---')


def main():
    with open('./books/frankenstein.txt', 'r') as f:
        file = f.read()

    print_report(file)


if __name__ == "__main__":
    main()
