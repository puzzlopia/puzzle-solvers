# puzzle-solvers
A Go puzzle solver used for solving and creating brain teasers. 

## Why another sliding blocks puzzle solver?
The intention is of course to be able to solve Sliding Blocks Puzzles, or SBPs.
But it is also to solve other puzzles, be it Sudokus, Sparse Sliding Blocks Puzzles (with lots of free space) or any kind.

I want [puzzlopia](http://www.puzzlopia.com) to have any type of puzzle, so I hope to add solvers for each type it deserves. And having a puzzle solver is a must if you also want to design new puzzles.

There are several inspiring attempts to develop solvers and finders. See [Neil's posts](https://nbickford.wordpress.com/2012/01/22/sliding-block-puzzles-part-3/) or [puzzlebeast](http://puzzlebeast.com/) for example. There is room for more ideas, algorithms and of course, puzzles.

Finally, a public repository can be viewed by other people, so they can just add more algorithms, better implementations, improve the solvers, fix bugs, etc.


## About this first implementation
I've selected Go language because it is similar to C++, equally powerful, and it addresses directly some interesting methodologies, as go-routines and parallelism (though I'm not using them, yet. This is my first Go project.)


## Sliding Block Puzzle Solver
The first solver implemented only works for SBPs. 
There is no restriction on the shape of pieces (yes, this is important for future puzzles).

As you can see if you run the code ('pennant_test' or 'solving-pennant'), the algorithm solves [Pennant](http://www.puzzlopia.com/puzzles/pennant/play) optimally (59 movs, less than a second on my laptop). I think it still requires adjustments in order to solve optimally any SBPs using the 'move metric'. Remember:
- 'Move metric': if you move around using only one piece (or block), it counts as one movement.
- 'Step metric': any movement of one position counts as one movement.

It is easy to implement an algorithm that finds solutions or optimal solutions using 'step metric'. [BFS](https://en.wikipedia.org/wiki/Depth-first_search) does pretty well. But everything changes when you use 'move metric', which seems to be the standard between SBPs experts. This is what current algorithm tries to do.

Currently, *the algorithm doesn't always find the optimal* (see 'checks' folder). Also it would be nice to add a _prune_ optimization.

## Using the solver
You only need to edit the 'main.go' file, uncomment 'solvingPennant()' and comment everything else. 

1. **Define the SBP puzzle**:
  You have to specify the initial state and the goal state. The structure is pretty straightforward: it consists of a matrix (by rows) of integers, 0 meaning free space and any positive value, a piece.
  Each piece must be marked with a different integer.
  To accelerate the algorithm, you can tell the algorithm that there are similar, interchangeable pieces: for example, the small 1x1 squares of
  Pennant. (They must have the same shape!).
  You can do it automatically:

  ```go
  myPuzzle.AutoAlikePieces()
  ```

  If there is a piece with the same shape as others and you need it not to be marked as 'alike' (this is common in a SBP where there are several 1x1 pieces and one of them is the piece used in the objective position), then you can specify that:

  ```go
  myPuzzle.SetNotAlikePiece(5)
  ```
  By default, all pieces are different.


2. **Adjust solver parameters**
  Basically, the limits:
  - 'MAX_DEPTH': max depth you enable the algorithm to reach. It should be equal or greater than the expected solution length in 'steps metric'.
  - 'MAX_STATES': you can limit the number of states explored, or you can leave it 0 and only limit max depth.


3. **Compile and run**. If the solver finds a solution, it will print it to the console. For Pennant example:

  ```bash
  Found! Path len:  59

  [id:4876] depth:83, GRID: [[7 6 3 3] [7 6 1 1] [0 0 5 4] [2 2 9 9] [2 2 8 8]]
  PATH<59>: [[4, 0, 1],[4, 0, 1],[5, 0, 1],[5, 0, 1],[2, 1, 0],[1, 0, -1],[1, 0, -1],[3, -1, 0],[5, -1, 0],[5, 0, 1],[2, 0, 1],[6, -1, 0],[6, -1, 0],[7, 0, -1],[8, 0, -1],[9, 0, -1],[4, 1, 0],[4, 1, 0],[5, 1, 0],[5, 1, 0],[2, 0, 1],[6, 0, 1],[7, -1, 0],[7, -1, 0],[8, 0, -1],[5, 0, -1],[4, -1, 0],[9, 0, 1],[8, 1, 0],[5, 0, -1],[5, 0, -1],[4, 0, -1],[4, 0, -1],[2, 1, 0],[3, 1, 0],[1, 0, 1],[1, 0, 1],[6, -1, 0],[4, -1, 0],[5, 0, 1],[7, 1, 0],[6, 0, -1],[4, -1, 0],[4, -1, 0],[5, -1, 0],[5, -1, 0],[7, 0, 1],[6, 1, 0],[6, 1, 0],[5, 0, -1],[5, -1, 0],[3, 0, -1],[3, 0, -1],[1, 1, 0],[4, 0, 1],[4, 0, 1],[5, 0, 1],[5, 0, 1],[3, -1, 0],[1, 0, -1],[4, 1, 0],[5, 0, 1],[3, 0, 1],[6, -1, 0],[6, -1, 0],[7, 0, -1],[2, 0, -1],[4, 1, 0],[4, 1, 0],[5, 1, 0],[5, 1, 0],[1, 0, 1],[3, 0, 1],[6, 0, 1],[7, -1, 0],[7, -1, 0],[2, 0, -1],[4, 0, -1],[4, -1, 0],[9, -1, 0],[8, 0, 1],[8, 0, 1],[2, 1, 0]]
  ```

## Dependencies
You will need the [fatih package](https://github.com/fatih/color), used to colorize the console output. Install:
```bash
go get github.com/fatih/color
```

## Next changes
Among the developments I would like to add:
- Refactor: I want to make cleaner interfaces, difficult because algorithms tend to perform better as more information is accessible, and this means larger interfaces and increased coupling.
- Improve current implementation: should be able to solve puzzlopia's [Ninja II](http://www.puzzlopia.com/puzzles/ninja-ii/play) (whithout waiting minutes).
- Prove (or at least know the limitations) that the algorithm always finds the optimal solution with 'move metric'.
- Add a user interface, probably a web interface.
- Add more features, like pieces with restrictions, puzzle and algorithm analytics, and more solvers (for sudokus, etc.)


## Contributors

- [Edgar Güeto](https://github.com/edgarweto)


