import ReactCodeMirror from "@uiw/react-codemirror";
import {CodeThemes} from "../../../app/utils";
import {useTheme} from "../../../provider/ThemeProvider";
import {PostgreSQL, sql} from "@codemirror/lang-sql";

type Props = {
    value: string,
    editable: boolean,
    onUpdate?: (value: string) => void,
}

export function QueryEditor(props: Props) {
    const {value, editable, onUpdate} = props
    const theme = useTheme();

    return (
        <ReactCodeMirror
            value={value}
            editable={editable}
            autoFocus={editable}
            minHeight={editable ? "80px" : "auto"}
            basicSetup={{lineNumbers: false, foldGutter: false, highlightActiveLine: false}}
            theme={CodeThemes[theme.mode]}
            extensions={[sql({dialect: PostgreSQL})]}
            onChange={onUpdate}
        />
    )
}
