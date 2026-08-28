import {Options} from "../../../core/widgets/options/Options"
import {SubContentBox} from "../../../shared/component/box/SubContentBox"
import {Options as ClusterOptions} from "../api/ClusterType"

type Props = {
    options: ClusterOptions,
    onUpdate: (options: ClusterOptions) => void,
}

// ClusterOptionsBox is the cluster options widget as a dialog shows it: folded
// away under its own heading, with the plugin selectors disabled, since inside
// a dialog the plugins come from the cluster list filter rather than the form.
export function ClusterOptionsBox(props: Props) {
    const {options, onUpdate} = props

    return (
        <SubContentBox label={"Cluster Options"} dense={true}>
            <Options options={options} onUpdate={onUpdate} disablePlugins={true}/>
        </SubContentBox>
    )
}
