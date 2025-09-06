import {QueryBoxCodeEditor} from "./QueryBoxCodeEditor";
import {InputBase, ToggleButton, ToggleButtonGroup, Tooltip} from "@mui/material";
import {QueryBoxInfo} from "./QueryBoxInfo";
import {DynamicInputs} from "../../view/input/DynamicInputs";
import {QueryVarietyOptions} from "../../../app/utils";
import {QueryRequest} from "../../../api/query/type";
import {SxPropsMap} from "../../../app/type";

const SX: SxPropsMap = {
    input: {fontSize: "inherit", padding: "0"},
    params: {fontSize: "inherit"},
    varieties: {lineHeight: "1.33", padding: "6px"},
}

type Props = {
    query: QueryRequest,
    onChange: (query: QueryRequest) => void,
}

export function QueryInfoEdit(props: Props) {
    const {query, onChange} = props

    return (
        <QueryBoxInfo
            type={query.type}
            editable={true}
            renderVarieties={renderVarietiesToggles()}
            renderDescription={renderDescription()}
            renderParams={renderParams()}
            renderQuery={renderQuery()}
        />
    )

    function renderDescription() {
        return (
            <InputBase
                sx={SX.input}
                fullWidth
                multiline
                placeholder={"Description"}
                value={query.description ?? ""}
                onChange={(e) => onChange({...query, description: e.target.value})}
            />
        )
    }

    function renderQuery() {
        return (
            <QueryBoxCodeEditor
                value={query.query}
                editable={true}
                autoFocus={false}
                onUpdate={(q) => onChange({...query, query: q})}
            />
        )
    }

    function renderParams() {
        return (
            <DynamicInputs
                inputs={query.params ?? []}
                placeholder={"Param: $"}
                editable={true}
                InputProps={SX.params}
                onChange={(p) => onChange({...query, params: p})}
            />
        )
    }

    function renderVarietiesToggles() {
        return (
            <ToggleButtonGroup
                size={"small"}
                value={query.varieties?.map(v => `${v}`)}
                onChange={(_e, v: string[]) => onChange({...query, varieties: v.map(s => parseInt(s))})}
            >
                {Object.entries(QueryVarietyOptions).map(([v, o]) => (
                    <ToggleButton key={v} sx={SX.varieties} value={v}>
                        <Tooltip title={o.label} placement={"top"}>
                            <span>{o.badge}</span>
                        </Tooltip>
                    </ToggleButton>
                ))}
            </ToggleButtonGroup>
        )
    }
}
