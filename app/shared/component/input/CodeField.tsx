import {useSnackbar} from "../../provider/SnackbarProvider"
import {TitleBox} from "../box/TitleBox"
import {CopyIconButton} from "../button/IconButtons"
import {CodeEditor} from "./CodeEditor"

type Props = {
    label: string,
    value: string,
    editable: boolean,
    hint?: string,
    placeholder?: string,
    minHeight?: string,
    copy?: boolean,
    onUpdate?: (value: string) => void,
}

// CodeField is a labelled block of code - the one way to show one, so a deploy
// command, a post script and a client's connection config are the same field at
// different sizes rather than things that happen to look similar.
export function CodeField(props: Props) {
    const {label, value, editable, hint, placeholder, minHeight = "80px", copy, onUpdate} = props
    const snackbar = useSnackbar()

    return (
        <TitleBox label={label} hint={hint} dense={true} collapsible={false} renderActions={renderCopy()}>
            <CodeEditor
                value={value}
                editable={editable}
                minHeight={minHeight}
                placeholder={placeholder}
                onUpdate={onUpdate}
            />
        </TitleBox>
    )

    function renderCopy() {
        if (!copy) return
        return <CopyIconButton tooltip={"Copy to your clipboard"} onClick={handleCopy}/>
    }

    function handleCopy() {
        navigator.clipboard.writeText(value).then(() => {
            snackbar("Copied to clipboard!", "info")
        })
    }
}
