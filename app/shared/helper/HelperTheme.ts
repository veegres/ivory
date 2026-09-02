import {ThemeOptions} from "@mui/material"

// NOTE: a field's height is its input line plus the vertical padding only - the
// outlined border lives on an absolutely positioned fieldset and adds nothing.
// LINE is MUI's own 1.4375em pinned to px, so a component that gives its input
// its own font (VaultInput's monospace) keeps the height of every other field.
const LINE = "23px"

const SIZE = {
    small: {pad: "1.5px"},
    medium: {pad: "4.5px"},
}

// NOTE: shrink repeats MUI's own value because the rest transform below
// outranks the variant it comes from, and would otherwise pin a floating
// label inside the field.
const LABEL = {
    small: "translate(14px, 1.5px) scale(1)",
    medium: "translate(14px, 4.5px) scale(1)",
    shrink: "translate(14px, -9px) scale(0.75)",
}

// NOTE: buttons carry their border in the box, so this is the whole height and
// the vertical padding has to go; it is a plain height so a component with its
// own toolbar sizing (Refresher, TriggerButton) still wins from sx.
const BUTTON = {
    small: {height: "26px", paddingTop: 0, paddingBottom: 0},
    medium: {height: "32px", paddingTop: 0, paddingBottom: 0},
}

const TAB = "36px"

const AUTOCOMPLETE = {
    small: {root: "0.5px", input: "1px", tag: "0px 3px", indicator: "0px"},
    medium: {root: "2px", input: "2.5px", tag: "1px 3px", indicator: "2px"},
}

export const ThemeComponents: ThemeOptions["components"] = {
    // NOTE: the :root variables are the app-wide size limits, reusable in any
    // sx as pure css strings, e.g. width: "var(--size-field)" or
    // flex: "1 1 min(var(--size-input), 100%)"
    MuiCssBaseline: {
        styleOverrides: {
            ":root": {
                "--size-field": "300px",
                "--size-input": "263px",
                "--size-tile": "350px",
                "--size-tile-height": "120px",
                "--size-dialog": "600px",
            },
        },
    },
    MuiOutlinedInput: {
        styleOverrides: {
            root: {
                ...getInputSize("medium"),
                "&.MuiInputBase-sizeSmall": getInputSize("small"),
            },
        },
    },
    MuiInputLabel: {
        styleOverrides: {
            root: {
                "&.MuiInputLabel-outlined": {
                    transform: LABEL.medium,
                    "&.MuiInputLabel-shrink": {transform: LABEL.shrink},
                    "&.MuiInputLabel-sizeSmall": {
                        transform: LABEL.small,
                        "&.MuiInputLabel-shrink": {transform: LABEL.shrink},
                    },
                },
            },
        },
    },
    MuiButton: {
        styleOverrides: {
            sizeSmall: BUTTON.small,
            sizeMedium: BUTTON.medium,
        },
    },
    MuiToggleButton: {
        styleOverrides: {
            sizeSmall: BUTTON.small,
            sizeMedium: BUTTON.medium,
        },
    },
    MuiTabs: {
        styleOverrides: {
            root: {minHeight: TAB},
        },
    },
    MuiTab: {
        styleOverrides: {
            root: {minHeight: TAB, paddingTop: 0, paddingBottom: 0},
        },
    },
    MuiAutocomplete: {
        styleOverrides: {
            root: {
                "& .MuiOutlinedInput-root": {
                    ...getAutocompleteSize("medium"),
                    "&.MuiInputBase-sizeSmall": getAutocompleteSize("small"),
                },
            },
        },
    },
}

type Size = keyof typeof SIZE

function getInputSize(size: Size) {
    const {pad} = SIZE[size]
    return {
        "& .MuiOutlinedInput-input": {height: LINE, paddingTop: pad, paddingBottom: pad},
        "&.MuiInputBase-multiline": {
            paddingTop: pad,
            paddingBottom: pad,
            "& .MuiOutlinedInput-input": {height: "auto", paddingTop: 0, paddingBottom: 0},
        },
    }
}

function getAutocompleteSize(size: Size) {
    const {root, input, tag, indicator} = AUTOCOMPLETE[size]
    return {
        paddingTop: root,
        paddingBottom: root,
        "& .MuiAutocomplete-input": {height: LINE, paddingTop: input, paddingBottom: input},
        "& .MuiAutocomplete-tag": {margin: tag},
        // MUI pins the adornment with a top offset computed for its own 28px
        // indicator, so a resized one has to be centred again.
        "& .MuiAutocomplete-endAdornment": {
            top: "50%",
            transform: "translateY(-50%)",
            "& .MuiIconButton-root": {padding: indicator},
        },
    }
}
