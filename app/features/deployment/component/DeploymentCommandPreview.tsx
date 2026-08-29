import {Box} from "@mui/material"

import {InfoColorBox} from "../../../shared/component/box/InfoColorBox"
import {InfoColorBoxRow} from "../../../shared/component/box/InfoColorBoxRow"
import {SubContentBox} from "../../../shared/component/box/SubContentBox"
import {CodeField} from "../../../shared/component/input/CodeField"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {DeployValues, interpolateCommand} from "../../../shared/helper/HelperUtils"
import {DeploymentPreviewNote} from "./DeploymentPreviewNote"

const SX: SxPropsMap = {
    preview: {display: "flex", flexDirection: "column", gap: 1},
}

type Props = {
    command: string,
    postScript?: string,
    // NOTE: whatever the screen already knows about the node - a value it has
    // no answer for yet leaves its placeholder standing
    values: DeployValues,
    defaultOpen?: boolean,
}

// DeploymentCommandPreview shows one node's command as it will run, with its
// post script below it. Both deploy screens read a command the same way: the
// cluster one folds it away because there is a form around it, the single-node
// one opens it because there is nothing else on the screen to read.
export function DeploymentCommandPreview(props: Props) {
    const {command, postScript, values, defaultOpen = false} = props

    return (
        /* NOTE: the badge rides on the toggle's own row - it says what is
           inside the section, so it belongs to the line that opens it */
        <SubContentBox label={"Preview"} renderActions={renderBadge()} defaultOpen={defaultOpen} dense={true}>
            {renderPreview()}
        </SubContentBox>
    )

    function renderBadge() {
        if (!postScript) return
        return (
            <InfoColorBoxRow>
                <InfoColorBox label={"post script"} title={"Runs inside the container once this node is up"}/>
            </InfoColorBoxRow>
        )
    }

    // NOTE: the hint sits above the code, not under it - it says how to read
    // what follows, which is no use once you have already read it
    function renderPreview() {
        return (
            <Box sx={SX.preview}>
                <DeploymentPreviewNote/>
                <CodeField
                    label={"Command"}
                    value={interpolateCommand(command, values)}
                    editable={false}
                    minHeight={"120px"}
                />
                {postScript && (
                    <CodeField
                        label={"Post Script"}
                        hint={"runs in the container once this node is up"}
                        value={interpolateCommand(postScript, values)}
                        editable={false}
                    />
                )}
            </Box>
        )
    }
}
