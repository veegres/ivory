import {PostgreSQL, sql} from "@codemirror/lang-sql"
import ReactCodeMirror from "@uiw/react-codemirror"

import {CodeThemes} from "../../../shared/helper/HelperUtils"
import {useSettings} from "../../../shared/provider/AppProvider"
import code from "../../../shared/style/codemirror.module.css"

type Props = {
    value: string,
    editable: boolean,
    autoFocus?: boolean,
    onUpdate?: (value: string) => void,
}

export function QueryBoxCodeEditor(props: Props) {
    const {value, editable, autoFocus, onUpdate} = props
    const settings = useSettings()

    return (
        <ReactCodeMirror
            className={code.simple}
            value={value}
            editable={editable}
            autoFocus={autoFocus}
            minHeight={"auto"}
            placeholder={"Query"}
            basicSetup={{lineNumbers: false, foldGutter: false, highlightActiveLine: false, highlightSelectionMatches: false}}
            theme={CodeThemes[settings.theme]}
            extensions={[sql({dialect: PostgreSQL})]}
            onChange={onUpdate}
        />
    )
}
