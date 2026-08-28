import {Decoration, DecorationSet, EditorView, MatchDecorator, ViewPlugin, ViewUpdate} from "@codemirror/view"
import {alpha, Box, Theme} from "@mui/material"
import ReactCodeMirror from "@uiw/react-codemirror"

import {SxPropsMap} from "../../helper/HelperType"
import {CodeThemes, PlaceholderPattern} from "../../helper/HelperUtils"
import {useSettings} from "../../provider/AppProvider"
import code from "../../style/codemirror.module.css"

const SX: SxPropsMap = {
    // NOTE: the padding goes on the scrolled content, not on this box. On the
    // box it stays put while the text scrolls underneath it, so a scrolled
    // command shows a line clipped in the padding above the first visible one.
    // The && doubles specificity so it wins over the codemirror module class.
    box: {
        bgcolor: "action.hover", borderRadius: 1, fontSize: "13px", overflow: "hidden",
        "&& .cm-content": {padding: "8px 0"},
        // NOTE: break-all, not the wrapping default - a docker run is one long
        // run of flags, paths and quoted script, and breaking only at spaces
        // left ragged half-empty lines wherever the next token was long
        "&& .cm-line": {padding: "0 10px", wordBreak: "break-all"},
        "&& .cm-scroller": {lineHeight: 1.5},
        // NOTE: warning, not an accent colour - an unresolved variable in a
        // preview is a field nobody filled in, and one in an editor is the
        // part of the command that is not yet real. Taken from the palette so
        // it holds up in either theme.
        [`&& .${code.placeholder}`]: {
            color: "warning.main",
            backgroundColor: (theme: Theme) => alpha(theme.palette.warning.main, 0.15),
        },
    },
}

type Props = {
    value: string,
    editable: boolean,
    placeholder?: string,
    autoFocus?: boolean,
    minHeight?: string,
    onUpdate?: (value: string) => void,
}

// NOTE: its own RegExp rather than the exported one - MatchDecorator keeps
// lastIndex on the object it is given, which would corrupt every other user
const placeholders = new MatchDecorator({
    regexp: new RegExp(PlaceholderPattern.source, "g"),
    decoration: Decoration.mark({class: code.placeholder}),
})

// PlaceholderHighlight picks the {{variables}} out of a command. In an editor
// they are what you are writing; in a preview one still showing is a field
// nobody has filled in, and either way they are easy to lose in a docker run.
const PlaceholderHighlight = ViewPlugin.fromClass(
    class {
        decorations: DecorationSet

        constructor(view: EditorView) {
            this.decorations = placeholders.createDeco(view)
        }

        update(update: ViewUpdate) {
            this.decorations = placeholders.updateDeco(update, this.decorations)
        }
    },
    {decorations: (plugin) => plugin.decorations},
)

// NOTE: wrapping rather than horizontal scrolling. A deploy command is one
// long line, and a scroller sized to its content puts the scrollbar directly
// under that line - in the middle of a box with a min height, not at its foot.
const Extensions = [EditorView.lineWrapping, PlaceholderHighlight]

export function CodeEditor(props: Props) {
    const {value, editable, placeholder, autoFocus, minHeight = "auto", onUpdate} = props
    const settings = useSettings()

    return (
        <Box sx={SX.box}>
            <ReactCodeMirror
                className={code.simple}
                value={value}
                editable={editable}
                autoFocus={autoFocus}
                minHeight={minHeight}
                placeholder={placeholder}
                basicSetup={{lineNumbers: false, foldGutter: false, highlightActiveLine: false, highlightSelectionMatches: false}}
                extensions={Extensions}
                theme={CodeThemes[settings.theme]}
                onChange={onUpdate}
            />
        </Box>
    )
}
