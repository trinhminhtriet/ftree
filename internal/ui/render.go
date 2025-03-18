package ui

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	"github.com/trinhminhtriet/ftree/internal/git"
	"github.com/trinhminhtriet/ftree/internal/state"
	t "github.com/trinhminhtriet/ftree/internal/tree"
	"github.com/trinhminhtriet/ftree/pkg/stack"
)

const (
	previewBytesLimit int64 = 10_000

	minHeight = 10
	minWidth  = 10

	arrow               = " <-"
	indentParent        = "│  "
	indentCurrent       = "├─ "
	indentCurrentLast   = "└─ "
	indentEmpty         = "   "
	emptydirContentName = "..."

	tooSmall                 = "too small =("
	binaryContentPlaceholder = "<binary content>"
	helpPreview              = "Press ? to toggle help"
)

type Renderer struct {
	Style       Stylesheet
	EdgePadding int
	offsetMem   int
	previewBuff [previewBytesLimit]byte
	gitStatus   map[string]git.GitStatus // Cache Git status
}

// NewRenderer initializes a Renderer with Git status.
func NewRenderer(style Stylesheet, edgePadding int) *Renderer {
	return &Renderer{
		Style:       style,
		EdgePadding: edgePadding,
		gitStatus:   make(map[string]git.GitStatus),
	}
}

func (r *Renderer) Render(s *state.State, winHeight, winWidth int) string {
	if winWidth < minWidth || winHeight < minHeight {
		return tooSmall
	}

	// Update Git status
	if r.gitStatus == nil || time.Since(time.Now()).Minutes() > 1 { // Refresh every minute
		newStatus, err := git.GetGitStatus(s.Tree.CurrentDir.Path)
		if err != nil {
			s.ErrBuf = fmt.Sprintf("Git status error: %v", err)
		} else {
			r.gitStatus = newStatus
		}
	}

	renderedHeading, headLen := r.renderHeading(s, winWidth)

	sectionWidth := int(math.Floor(0.5 * float64(winWidth)))

	renderedTree := r.renderTree(s.Tree, winHeight-headLen, sectionWidth)

	var rightPane string
	if s.HelpToggle {
		renderedHelp, helpLen := r.renderHelp(sectionWidth)
		renderedContent := r.renderSelectedFileContent(s.Tree, winHeight-headLen-helpLen, sectionWidth)
		rightPane = lipgloss.JoinVertical(lipgloss.Left, renderedHelp, renderedContent)
	} else {
		renderedContent := r.renderSelectedFileContent(s.Tree, winHeight-headLen, sectionWidth)
		rightPane = renderedContent
	}

	renderedTreeWithContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		renderedTree,
		rightPane,
	)

	return renderedHeading + "\n" + renderedTreeWithContent
}

func (r *Renderer) renderHeading(s *state.State, width int) (string, int) {
	selected := s.Tree.GetSelectedChild()

	path := s.Tree.CurrentDir.Path + "/..."
	changeTime := "--"
	size := "0 B"
	perm := "--"
	gitStatus := ""

	if selected != nil {
		path = selected.Path
		changeTime = selected.Info.ModTime().Format(time.RFC822)
		size = formatSize(float64(selected.Info.Size()), 1024.0)
		perm = selected.Info.Mode().String()
		if status, ok := r.gitStatus[selected.Path]; ok {
			if status.Staged {
				gitStatus = " [Staged]"
			} else if status.Modified {
				gitStatus = " [Modified]"
			} else if status.Untracked {
				gitStatus = " [Untracked]"
			}
		}
	}

	markedPath := ""
	if s.Tree.Marked != nil {
		markedPath = s.Tree.Marked.Path
	}

	operationBar := fmt.Sprintf(": %s", s.OpBuf.Repr())
	if markedPath != "" {
		operationBar += fmt.Sprintf(" [%s]", markedPath)
	}

	if s.OpBuf.IsInput() {
		operationBar += fmt.Sprintf(" │ %s │", r.Style.OperationBarInput.Render(string(s.InputBuf)))
	}

	rawPath := "> " + path + gitStatus

	finfo := fmt.Sprintf(
		"%s %s %v %s %s",
		r.Style.FinfoPermissions.Render(perm),
		r.Style.FinfoSep.Render("│"),
		r.Style.FinfoLastUpdated.Render(changeTime),
		r.Style.FinfoSep.Render("│"),
		r.Style.FinfoSize.Render(size),
	)

	header := []string{
		r.Style.SelectedPath.Render(rawPath) +
			strings.Repeat(
				" ",
				max(width-utf8.RuneCountInString(rawPath)-utf8.RuneCountInString(helpPreview), 0),
			) +
			r.Style.HelpMsg.Render(helpPreview),
		finfo,
		r.Style.OperationBar.Render(operationBar),
		r.Style.ErrBar.Render(s.ErrBuf),
	}
	return strings.Join(header, "\n"), len(header)
}

func (r *Renderer) renderHelp(width int) (string, int) {
	help := []string{
		"j / arr down   Select next child",
		"k / arr up     Select previous child",
		"h / arr left   Move up a dir",
		"l / arr right  Enter selected directory",
		"if / id        Create file (if) / directory (id)",
		"d              Move selected child (then 'p' to paste)",
		"y              Copy selected child (then 'p' to paste)",
		"D              Delete selected child",
		"r              Rename selected child",
		"e              Edit selected file in $EDITOR",
		"gg             Go to top child",
		"G              Go to last child",
		"enter          Collapse / expand directory",
		"esc            Clear error / stop operation",
		"q / ctrl+c     Exit",
		"g?             Toggle Git status details", // New help item
	}
	return r.Style.
		HelpContent.
		MaxWidth(width).
		MarginRight(width).
		Render(strings.Join(help, "\n")), len(help) + 1
}

func (r *Renderer) renderTree(tree *t.Tree, height, width int) string {
	renderedTreeLines, selectedRow := r.renderTreeFull(tree, width)
	croppedTreeLines := r.cropTree(renderedTreeLines, selectedRow, height)

	// Add Git status to tree lines if enabled (mock toggle with g?)
	if tree.CurrentDir.Path != "" && r.gitStatus != nil {
		for i, line := range croppedTreeLines {
			for path, status := range r.gitStatus {
				if strings.Contains(line, path) {
					var statusStr string
					if status.Staged {
						statusStr = " [S]"
					} else if status.Modified {
						statusStr = " [M]"
					} else if status.Untracked {
						statusStr = " [U]"
					}
					croppedTreeLines[i] += r.Style.TreeGitStatus.Render(statusStr)
					break
				}
			}
		}
	}

	treeStyle := lipgloss.
		NewStyle().
		MaxWidth(width).
		MarginRight(width)

	return treeStyle.Render(strings.Join(croppedTreeLines, "\n"))
}

func (r *Renderer) renderSelectedFileContent(tree *t.Tree, height, width int) string {
	n, err := tree.ReadSelectedChildContent(r.previewBuff[:], previewBytesLimit)
	if err != nil {
		return r.Style.ContentPreview.Render(fmt.Sprintf("Error: %v", err))
	}
	content := r.previewBuff[:n]

	contentStyle := r.Style.ContentPreview.MaxWidth(width - 1)
	var contentLines []string
	if !utf8.Valid(content) {
		contentLines = []string{binaryContentPlaceholder}
	} else {
		contentLines = strings.Split(string(content), "\n")
		contentLines = contentLines[:min(max(height, 0), len(contentLines))]
	}
	return contentStyle.Render(strings.Join(contentLines, "\n"))
}

func (r *Renderer) cropTree(lines []string, currentLine int, height int) []string {
	linesLen := len(lines)

	offset := r.offsetMem
	limit := linesLen

	if currentLine+1 > height+offset-r.EdgePadding {
		offset = max(min(currentLine+1-height+r.EdgePadding, linesLen-height), 0)
	}
	if currentLine < r.EdgePadding+offset {
		offset = max(currentLine-r.EdgePadding, 0)
	}
	r.offsetMem = offset
	limit = min(height+offset, linesLen)
	return lines[offset:limit]
}

func (r *Renderer) renderTreeFull(tree *t.Tree, width int) ([]string, int) {
	linen := -1
	currentLine := 0

	type stackEl struct {
		*t.Node
		string
		bool
	}
	lines := []string{}
	s := stack.NewStack(stackEl{tree.Root, "", false})

	for s.Len() > 0 {
		el := s.Pop()
		linen += 1

		node := el.Node
		isLast := el.bool
		parentIndent := el.string

		var indent string
		if node == tree.Root {
			indent = ""
		} else if isLast {
			indent = parentIndent + indentCurrentLast
			parentIndent = parentIndent + indentEmpty
		} else {
			indent = parentIndent + indentCurrent
			parentIndent = parentIndent + indentParent
		}

		if node == nil {
			continue
		}

		name := node.Info.Name()
		nameRuneCountNoStyle := utf8.RuneCountInString(name)
		indentRuneCount := utf8.RuneCountInString(indent)

		if nameRuneCountNoStyle+indentRuneCount > width-6 {
			name = string([]rune(name)[:max(0, width-indentRuneCount-6)]) + "..."
		}

		indent = r.Style.TreeIndent.Render(indent)

		if node.Info.IsDir() {
			name = r.Style.TreeDirectoryName.Render(name) // Fixed typo: Direcotry -> Directory
		} else if node.Info.Mode()&os.ModeSymlink == os.ModeSymlink {
			name = r.Style.TreeLinkName.Render(name)
		} else {
			name = r.Style.TreeRegularFileName.Render(name)
		}

		if tree.Marked == node {
			name = r.Style.TreeMarkedNode.Render(name)
		}

		repr := indent + name

		if tree.GetSelectedChild() == node {
			repr += r.Style.TreeSelectionArrow.Render(arrow)
			currentLine = linen
		}
		lines = append(lines, repr)

		if node.Children != nil {
			if len(node.Children) == 0 && tree.CurrentDir == node {
				emptyIndent := r.Style.TreeIndent.Render(parentIndent + indentCurrentLast)
				lines = append(lines, emptyIndent+emptydirContentName+r.Style.TreeSelectionArrow.Render(arrow))
				currentLine = linen + 1
			}
			for i := len(node.Children) - 1; i >= 0; i-- {
				ch := node.Children[i]
				s.Push(stackEl{ch, parentIndent, i == len(node.Children)-1})
			}
		}
	}
	return lines, currentLine
}

var sizes = [...]string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}

func formatSize(s float64, base float64) string {
	unitsLimit := len(sizes)
	i := 0
	for s >= base && i < unitsLimit-1 { // -1 to avoid out-of-bounds
		s = s / base
		i++
	}
	f := "%.0f %s"
	if i > 1 {
		f = "%.2f %s"
	}
	return fmt.Sprintf(f, s, sizes[i])
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
