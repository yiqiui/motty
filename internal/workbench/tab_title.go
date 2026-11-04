package workbench

import (
	"strings"

	"github.com/vivaui/viva/ui"
	layoutcore "motty/internal/workspace"
)

// tabLabel is the name a tab shows. A tab follows the title its active
// terminal reports (OSC 0 and 2, which tmux sets from its window and session
// names with set-titles) when the setting is on, unless the user named the tab
// themself. The stored workspace title stays the host's or the user's name: the
// reported one is a property of what is running right now, so it is neither
// saved nor used to address the session in questions.
func (state *State) tabLabel(workspace layoutcore.WorkspaceInfo) string {
	if !state.terminalTabTitle || workspace.Named {
		return workspace.Title
	}
	session := state.tabSession(workspace)
	if session == nil || session.term == nil || session.term.screen == nil {
		return workspace.Title
	}
	if reported := strings.TrimSpace(session.term.screen.Title()); reported != "" {
		return reported
	}
	return workspace.Title
}

// tabSession returns the session behind the pane the tab snapshot currently
// represents. Tab labels and badges share this lookup so a pane activation
// updates both together.
func (state *State) tabSession(workspace layoutcore.WorkspaceInfo) *Session {
	pane, ok := state.deck.Pane(workspace.ID, workspace.ActivePane)
	if !ok || pane.SessionID == 0 {
		return nil
	}
	session, _ := state.sessionByIDPointer(uint32(pane.SessionID))
	return session
}

// tabLabelFits reports whether label fits the tab's text area at the size and
// weight the tab draws its title. declareWorkspaceTab reserves the same fixed
// chrome a painted tab occupies — the badge, host light and close button plus
// their margins and the row's inner padding — and lets the label take what
// remains, exactly as the visible title does, so the prediction matches what
// the eye sees (issue #352).
func (state *State) tabLabelFits(canvas *ui.Canvas, workspace layoutcore.WorkspaceInfo, label string, width float32, compact bool) bool {
	const (
		padLeft  = 4
		badge    = 20 + 7
		light    = 9 + 7
		close    = 20 + 7
		reserve  = padLeft + badge + light + close
	)
	// A compact tab shows only a host-color square and no label, so it neither
	// elides nor needs a hover hint.
	if compact {
		return true
	}
	avail := width - reserve
	if avail <= 0 {
		return false
	}
	return canvas.TextWidthWeight(label, metrics.FontSizeS, ui.WeightMedium) <= avail
}
