import {SxPropsMap} from "../../../type/general";
import {QueryBoxCodeEditor} from "./QueryBoxCodeEditor";
import {Query} from "../../../type/query";
import {QueryBoxInfo} from "./QueryBoxInfo";
import {QueryVarieties} from "./QueryVarieties";
import {Avatar, Box, Chip} from "@mui/material";

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", gap: 1},
    paper: {padding: "10px", background: "rgba(145,145,145,0.1)", borderRadius: "10px"},
    varieties: {display: "grid", gridAutoColumns: "minmax(0, 1fr)", gridAutoFlow: "column", gap: 2},
    params: {display: "flex", gap: 1, flexWrap: "wrap"},
    no: {display: "flex", alignItems: "center", justifyContent: "center", height: "100%", color: "text.secondary"},
}

type Props = {
    query: Query,
}

export function QueryInfoView(props: Props) {
    const {type, params, description, varieties, custom} = props.query

    return (
        <QueryBoxInfo
            type={type}
            editable={false}
            renderVarieties={renderVarieties()}
            renderDescription={renderDescription()}
            renderParams={renderParams()}
            renderQuery={renderQuery()}
        />
    )

    function renderDescription() {
        if (!description) return renderNoElement("DESCRIPTION")

        return description
    }

    function renderVarieties() {
        if (!varieties) return renderNoElement("LABELS")

        return (
            <QueryVarieties varieties={varieties}/>
        )
    }

    function renderParams() {
        if (!params) return renderNoElement("PREPARED PARAMS")
        return (
            <Box sx={SX.params}>
                {params.map((input, index) => (
                    <Chip
                        key={index}
                        avatar={<Avatar>${index + 1}</Avatar>}
                        label={input}
                        disabled={!input}
                    />
                ))}
            </Box>
        )
    }

    function renderQuery() {
        return (
            <QueryBoxCodeEditor value={custom} editable={false}/>
        )
    }

    function renderNoElement(name: string) {
        return (
            <Box sx={SX.no}>NO {name}</Box>
        )
    }
}
