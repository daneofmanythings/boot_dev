import tkinter as tk
from maze import Maze
from config import WINDOWHIEGHT, WINDOWWIDTH


class Window:
    def __init__(self, width, height):
        self._root = tk.Tk()
        self._root.title = "title"
        self._root.protocol("WM_DELETE_WINDOW", self.close)

        self._canvas = tk.Canvas(height=height, width=width)
        self._canvas.pack()

        self._is_running = False

    def draw_line(self, line, fill_color):
        line.draw(self._canvas, fill_color)

    def redraw(self):
        self._root.update_idletasks()
        self._root.update()

    def wait_for_close(self):
        self._is_running = True
        while self._is_running:
            self.redraw()

    def close(self):
        self._is_running = False


def main():
    win = Window(WINDOWWIDTH, WINDOWHIEGHT)
    maze = Maze(
        x1=20,
        y1=20,
        num_rows=40,
        num_cols=40,
        size_x=15,
        size_y=15,
        win=win,
        # seed=10
    )

    maze.solve()

    win.wait_for_close()


if __name__ == "__main__":
    main()
