import {TitleBox} from "../box/TitleBox"
import {CodeEditor} from "./CodeEditor"

type Props = {
    label: string,
    value: string,
    editable: boolean,
    hint?: string,
    placeholder?: string,
    minHeight?: string,
    onUpdate?: (value: string) => void,
}

// CodeField is a labelled block of code - the deploy dialog's only way to show
// one, so a command and a post script are the same field at different sizes
// rather than two things that happen to look similar.
export function CodeField(props: Props) {
    const {label, value, editable, hint, placeholder, minHeight = "80px", onUpdate} = props

    return (
        <TitleBox label={label} hint={hint} dense={true} collapsible={false}>
            <CodeEditor
                value={value}
                editable={editable}
                minHeight={minHeight}
                placeholder={placeholder}
                onUpdate={onUpdate}
            />
        </TitleBox>
    )
}
