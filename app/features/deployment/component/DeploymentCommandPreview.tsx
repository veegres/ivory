import {Box} from "@mui/material"

import {InfoColorBox} from "../../../shared/component/box/InfoColorBox"
import {InfoColorBoxRow} from "../../../shared/component/box/InfoColorBoxRow"
import {TitleBox} from "../../../shared/component/box/TitleBox"
import {CodeField} from "../../../shared/component/input/CodeField"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {DeployPasswordMask, DeployValues, interpolateCommand} from "../../../shared/helper/HelperUtils"

const SX: SxPropsMap = {
    preview: {display: "flex", flexDirection: "column", gap: 1},
}

type Props = {
    command: string,
    postScripts?: string[],
    // NOTE: whatever the screen already knows about the node - a value it has
    // no answer for yet leaves its placeholder standing
    values: DeployValues,
    defaultOpen?: boolean,
}

// DeploymentCommandPreview shows one node's command as it will run, with its
// post script below it. Both deploy screens read a command the same way: the
// cluster one folds it away because there is a form around it, the single-node
// one opens it because there is nothing else on the screen to read. Everything
// is real except the password, so the mask is the only thing worth explaining -
// the section's own hint does that, since a reader has to know the deployed
// command is not quite the one on screen.
export function DeploymentCommandPreview(props: Props) {
    const {command, postScripts, values, defaultOpen = false} = props

    return (
        /* NOTE: the badge rides on the toggle's own row - it says what is
           inside the section, so it belongs to the line that opens it */
        <TitleBox
            label={"Preview"}
            hint={`This is what will run on the node, with ${DeployPasswordMask} standing in for the password: the server substitutes the real one from the database credentials above, so it never reaches the browser. It does reach the node, though - docker keeps the real command it ran, so anyone with docker access there can read the password back out with docker inspect.`}
            renderActions={renderBadge()}
            defaultOpen={defaultOpen}
            dense={true}
        >
            {renderPreview()}
        </TitleBox>
    )

    function renderBadge() {
        if (!postScripts?.length) return
        return (
            <InfoColorBoxRow>
                <InfoColorBox label={"post script"} title={"Runs inside the container once this node is up"}/>
            </InfoColorBoxRow>
        )
    }

    function renderPreview() {
        return (
            <Box sx={SX.preview}>
                <CodeField
                    label={"Command"}
                    value={interpolateCommand(command, values)}
                    editable={false}
                    minHeight={"120px"}
                />
                {!!postScripts?.length && (
                    <CodeField
                        label={"Post Script"}
                        hint={"each line runs in the container once this node is up"}
                        value={postScripts.map((s) => interpolateCommand(s, values)).join("\n")}
                        editable={false}
                    />
                )}
            </Box>
        )
    }
}
