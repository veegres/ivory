import ReactCodeMirror from "@uiw/react-codemirror";
import {CodeThemes} from "../../../app/utils";
import {useAppearance} from "../../../provider/AppearanceProvider";
import {PostgreSQL, sql} from "@codemirror/lang-sql";

type Props = {
    value: string,
    editable: boolean,
    onUpdate?: (value: string) => void,
}

export function QueryEditor(props: Props) {
    const {value, editable, onUpdate} = props
    const appearance = useAppearance();

    return (
        <ReactCodeMirror
            value={value}
            editable={editable}
            autoFocus={editable}
            minHeight={editable ? "80px" : "auto"}
            placeholder={"Query"}
            basicSetup={{lineNumbers: false, foldGutter: false, highlightActiveLine: false, highlightSelectionMatches: false}}
            theme={CodeThemes[appearance.theme]}
            extensions={[sql({dialect: PostgreSQL})]}
            onChange={onUpdate}
        />
    )
}
