import {useQuery, useQueryClient} from "@tanstack/react-query";
import {instanceApi} from "../app/api";
import {useEffect, useMemo, useRef} from "react";
import {combineInstances, createInstanceColors, getDomain, initialInstance, isSidecarEqual} from "../app/utils";
import {useStore, useStoreAction} from "../provider/StoreProvider";
import {Cluster, DetectionType, InstanceDetection} from "../type/cluster";
import {Sidecar} from "../type/common";
import {InstanceMap} from "../type/instance";

export function useInstanceDetection(cluster: Cluster, instances: Sidecar[]): InstanceDetection {
    const {activeCluster} = useStore()
    const {setCluster, setClusterInfo, setWarnings} = useStoreAction()
    const isClusterActive = !!activeCluster && cluster.name === activeCluster.cluster.name

    const defaultDetection = useRef<DetectionType>("auto")
    const defaultSidecar = useRef(instances[0])
    const previousData = useRef<InstanceMap>({})

    const queryKey = useMemo(() => ["instance/overview", cluster.name], [cluster.name])
    const queryClient = useQueryClient();
    const query = useQuery({
        queryKey: ["instance/overview", cluster.name],
        queryFn: () => instanceApi.overview({
            ...defaultSidecar.current,
            credentialId: cluster.credentials.patroniId,
            certs: cluster.certs
        }),
        retry: false,
    })
    const {errorUpdateCount, refetch, data, isFetching} = query
    const instanceMap = useMemo(handleMemoInstanceMap, [data])
    previousData.current = instanceMap

    const combine = useMemo(handleMemoCombine, [instances, instanceMap])
    const colors =  useMemo(handleMemoColors, [combine.combinedInstanceMap])
    const defaultInstance = handleDefaultInstance()

    useEffect(handleEffectNextRequest, [errorUpdateCount, instances, refetch])
    useEffect(handleEffectDetectionChange, [isClusterActive, activeCluster, instances, refetch, queryClient, queryKey])

    useEffect(handleEffectRequestUpdate, [isClusterActive, cluster, defaultInstance, combine, setCluster, setClusterInfo])
    useEffect(handleEffectWarningsUpdate, [cluster.name, combine.warning, setWarnings])

    return {
        defaultInstance,
        combinedInstanceMap: combine.combinedInstanceMap,
        warning: combine.warning,
        detection: isClusterActive ? activeCluster.detection : defaultDetection.current,
        colors,
        active: isClusterActive,
        fetching: isFetching,
        refetch: handleRefetch,
    }

    /**
     * When we refetch in `auto` detection we need to start from scratch, that is why we change
     * `defaultSidecar` value to the first one. After that we clean `useQuery`
     * state and start fetching from scratch with error count equal 0. That
     * is why `handleEffectNextRequest` continue working. We don't send request additional
     * request to the master node.
     */
    async function handleRefetch() {
        if (defaultDetection.current === "auto") {
            defaultSidecar.current = instances[0]
        }
        queryClient.removeQueries({queryKey})
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
        if (isClusterActive) {
            const oldDetection = defaultDetection.current
            const newDetection = activeCluster.detection
            const isSidecarNotEqual = !isSidecarEqual(defaultSidecar.current, activeCluster.defaultInstance.sidecar)

            if (newDetection === "manual" && (isSidecarNotEqual || oldDetection === "auto")) {
                defaultSidecar.current = activeCluster.defaultInstance.sidecar
                defaultDetection.current = activeCluster.detection
                queryClient.removeQueries({queryKey})
                refetch().then()
            }

            if (newDetection === "auto" && oldDetection === "manual") {
                defaultSidecar.current = instances[0]
                defaultDetection.current = activeCluster.detection
                queryClient.removeQueries({queryKey})
                refetch().then()
            }
        }
    }

    /**
     * Update store each time when request is happened
     */
    function handleEffectRequestUpdate() {
        if (isClusterActive) {
            const detection = defaultDetection.current
            setCluster({detection, cluster, defaultInstance, ...combine})
        }

        return () => {
            if (isClusterActive) setCluster(undefined)
        }
    }

    function handleEffectWarningsUpdate() {
        setWarnings(cluster.name, combine.warning)
        return () => {
            if (combine.warning) setWarnings(cluster.name, false)
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
        return createInstanceColors(combine.combinedInstanceMap)
    }

    function handleMemoCombine() {
        return combineInstances(instances, instanceMap)
    }

    /**
     * Either find leader or set query that we send request to.
     * P.S. we cannot use memo for this function because `combine` doesn't change,
     * but we still need to change it every time for clusters without any success request
     */
    function handleDefaultInstance() {
        const map = combine.combinedInstanceMap

        if (defaultDetection.current === "manual") {
            if (isClusterActive) {
                const sidecar = activeCluster.defaultInstance.sidecar
                return map[getDomain(sidecar)] ?? initialInstance(sidecar)
            } else {
                return map[getDomain(defaultSidecar.current)] ?? initialInstance(defaultSidecar.current)
            }
        } else {
            const leader = Object.values(map).find(i => i.leader)
            return leader ?? map[getDomain(defaultSidecar.current)] ?? initialInstance(defaultSidecar.current)
        }
    }
}
