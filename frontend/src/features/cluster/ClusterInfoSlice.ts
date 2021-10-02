import {createAsyncThunk, createSlice} from '@reduxjs/toolkit';
import {RootState} from '../../app/store';
import {restClient} from "../../app/hooks";

export interface Node {
    name: string,
    timeline: number,
    lag: number,
    state: string,
    host: string,
    role: string,
    port: number,
    api_url: string
}

export interface ClusterInfoState {
    data: Node[];
    error?: string;
    status: 'idle' | 'loading' | 'failed';
}

const initialState: ClusterInfoState = {
    data: [],
    status: 'idle',
};

export const incrementAsync = createAsyncThunk(
    'cluster/fetchInfo',
    async (node: string, { rejectWithValue }) => {
        try {
            return await restClient<{ members: Node[] }>(`/node/${node}/cluster`)
        } catch (error) {
            return rejectWithValue(error.response.data)
        }
    }
);

export const clusterInfoSlice = createSlice({
    name: 'cluster/info',
    initialState,
    reducers: {},
    extraReducers: (builder) => {
        builder
            .addCase(incrementAsync.pending, (state) => {
                state.status = 'loading';
            })
            .addCase(incrementAsync.fulfilled, (state, action) => {
                state.status = 'idle';
                state.data = action.payload.response.members;
            })
            .addCase(incrementAsync.rejected, (state, action) => {
                state.status = 'idle';
                state.error = action.error.message;
            });
    },
});

export const selectClusterData = (state: RootState) => state.clusterInfo.data;
