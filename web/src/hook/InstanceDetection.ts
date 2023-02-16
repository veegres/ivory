import {DetectionType, InstanceDetection, InstanceMap, Sidecar} from "../app/types";
import {useQuery} from "@tanstack/react-query";
import {instanceApi} from "../app/api";
import {useEffect, useMemo, useRef} from "react";
import {combineInstances, createInstanceColors, getDomain, initialInstance, isSidecarEqual} from "../app/utils";
import {useStore} from "../provider/StoreProvider";

export function useInstanceDetection(name: string, instances: Sidecar[]): InstanceDetection {
    const {store: {activeCluster}, setCluster} = useStore()
    const isActive = !!activeCluster && name === activeCluster.cluster.name

    const defaultDetection = useRef<DetectionType>("auto")
    const defaultSidecar = useRef(instances[0])
    const previousData = useRef<InstanceMap>({})

    const query = useQuery(
        ["instance/overview", name],
        () => instanceApi.overview({...defaultSidecar.current, cluster: name}),
        {retry: false}
    )
    const {errorUpdateCount, refetch, remove, data, isFetching} = query
    const instanceMap = useMemo(handleMemoInstanceMap, [data])
    previousData.current = instanceMap

    const colors =  useMemo(handleMemoColors, [instanceMap])
    const combine = useMemo(handleMemoCombine, [instances, instanceMap])
    const defaultInstance = useMemo(handleMemoDefaultInstance, [activeCluster, combine])

    useEffect(handleEffectNextRequest, [errorUpdateCount, instances, refetch])
    useEffect(handleEffectDetectionChange, [isActive, activeCluster, instances, refetch, remove])

    // we ignore this line cause this effect uses activeCluster and setCluster
    // which are always changing in this function, and it causes endless recursion
    // eslint-disable-next-line
    useEffect(handleRequestUpdate, [isActive, defaultInstance, combine])

    return {
        defaultInstance,
        combinedInstanceMap: combine.combinedInstanceMap,
        warning: combine.warning,
        colors,
        fetching: isFetching,
        refetch: handleRefetch,
    }

    /**
     * When we refetch we need to start from scratch, that is why we change
     * `defaultSidecar` value to the first one. After that we clean `useQuery`
     * state and start fetching from scratch with error count equal 0. That
     * is why `handleEffectNextRequest` continue working.
     */
    async function handleRefetch() {
        if (defaultDetection.current === "auto") {
            defaultSidecar.current = instances[0]
        }
        remove()
        await refetch()
    }

    /**
     * We do refetch in each error, we shouldn't refetch when:
     * - error count is bigger than our amount of instances
     * - sidecar has default value
     * - we already fetched this instance, and it still the same
     * - detection is manual
     */
    function handleEffectNextRequest() {
        if (errorUpdateCount < instances.length &&
            defaultSidecar.current !== undefined &&
            defaultSidecar.current !== instances[errorUpdateCount] &&
            defaultDetection.current !== "manual"
        ) {
            defaultSidecar.current = instances[errorUpdateCount]
            refetch().then()
        }
    }

    /**
     * When we change detection we want to refresh request
     * - manual - refresh each time when instance changed
     * - auto   - refresh from scratch if detection was before manual
     */
    function handleEffectDetectionChange() {
        if (isActive) {
            const oldDetection = defaultDetection.current
            const newDetection = activeCluster.detection
            const isSidecarNotEqual = !isSidecarEqual(defaultSidecar.current, activeCluster.defaultInstance.sidecar)

            if (newDetection === "manual" && isSidecarNotEqual) {
                defaultSidecar.current = activeCluster.defaultInstance.sidecar
                defaultDetection.current = activeCluster.detection
                remove()
                refetch().then()
            }

            if (newDetection === "auto" && oldDetection === "manual") {
                defaultSidecar.current = instances[0]
                defaultDetection.current = activeCluster.detection
                remove()
                refetch().then()
            }
        }
    }

    /**
     * Update store each time when request is happened
     */
    function handleRequestUpdate() {
        if (isActive) {
            setCluster({...activeCluster, defaultInstance, ...combine})
        }
    }

    function handleMemoInstanceMap() {
        if (!data && defaultDetection.current === "auto") {
            return previousData.current
        } else {
            return data ?? {}
        }
    }

    function handleMemoColors() {
        return createInstanceColors(instanceMap)
    }

    function handleMemoCombine() {
        return combineInstances(instances, instanceMap)
    }

    /**
     * Either find leader or set query that we send request to.
     */
    function handleMemoDefaultInstance() {
        const map = combine.combinedInstanceMap

        if (activeCluster && defaultDetection.current === "manual") {
            const selected = map[getDomain(activeCluster.defaultInstance.sidecar)]
            return selected ?? initialInstance(activeCluster.defaultInstance.sidecar)
        } else {
            const leader = Object.values(map).find(i => i.leader)
            return leader ?? map[getDomain(defaultSidecar.current)]
        }
    }
}
