import ReactCodeMirror from "@uiw/react-codemirror";
import {CodeThemes} from "../../../app/utils";
import {useAppearance} from "../../../provider/AppearanceProvider";
import {PostgreSQL, sql} from "@codemirror/lang-sql";
import code from "../../../style/codemirror.module.css";

type Props = {
    value: string,
    editable: boolean,
    autoFocus?: boolean,
    onUpdate?: (value: string) => void,
}

export function QueryBoxCodeEditor(props: Props) {
    const {value, editable, autoFocus, onUpdate} = props
    const appearance = useAppearance();

    return (
        <ReactCodeMirror
            className={code.simple}
            value={value}
            editable={editable}
            autoFocus={autoFocus}
            minHeight={"auto"}
            placeholder={"Query"}
            basicSetup={{lineNumbers: false, foldGutter: false, highlightActiveLine: false, highlightSelectionMatches: false}}
            theme={CodeThemes[appearance.theme]}
            extensions={[sql({dialect: PostgreSQL})]}
            onChange={onUpdate}
        />
    )
}
