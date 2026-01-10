<template>
    <el-card class="task-queue-panel" shadow="hover">
        <template #header>
            <div class="task-queue-header">
                <span><el-icon>
                        <View />
                    </el-icon> 任务队列状态</span>
                <div>
                    <el-button type="text" @click="fetchQueueStatus"><el-icon>
                            <Refresh />
                        </el-icon> 刷新</el-button>
                </div>
            </div>
        </template>

        <el-row :gutter="20" class="task-queue-statistics">
            <el-col :span="6">
                <el-statistic title="向量任务 - 进行中" :value="embeddingStats.active">
                    <template #suffix>个</template>
                </el-statistic>
            </el-col>
            <el-col :span="6">
                <el-statistic title="向量任务 - 总计" :value="embeddingStats.total">
                    <template #suffix>个</template>
                </el-statistic>
            </el-col>
            <el-col :span="6">
                <el-statistic title="摘要任务 - 进行中" :value="summaryStats.active">
                    <template #suffix>个</template>
                </el-statistic>
            </el-col>
            <el-col :span="6">
                <el-statistic title="摘要任务 - 总计" :value="summaryStats.total">
                    <template #suffix>个</template>
                </el-statistic>
            </el-col>
        </el-row>

        <el-tabs type="border-card">
            <el-tab-pane label="向量生成任务">
                <template #label>
                    <span>
                        <el-icon>
                            <DataAnalysis />
                        </el-icon>
                        向量生成
                        <el-badge
                            v-if="queueStatus.pending_embedding_tasks.length + queueStatus.processing_embedding_tasks.length > 0"
                            :value="queueStatus.pending_embedding_tasks.length + queueStatus.processing_embedding_tasks.length"
                            class="item" />
                    </span>
                </template>

                <div v-if="queueStatus.processing_embedding_tasks.length > 0" class="task-section">
                    <el-divider content-position="left" class="task-divider">
                        <el-tag type="success" size="small">正在处理 ({{ queueStatus.processing_embedding_tasks.length
                        }})</el-tag>
                    </el-divider>
                    <el-table :data="queueStatus.processing_embedding_tasks" size="small" stripe>
                        <el-table-column prop="opera_id" label="视频ID" width="100" />
                        <el-table-column prop="id" label="任务ID" show-overflow-tooltip />
                        <el-table-column prop="created_at" label="创建时间" width="180">
                            <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
                        </el-table-column>
                        <el-table-column prop="updated_at" label="更新时间" width="180">
                            <template #default="{ row }">{{ formatTime(row.updated_at) }}</template>
                        </el-table-column>
                        <el-table-column label="状态" width="100">
                            <template #default>
                                <el-tag type="success" size="small">
                                    <el-icon class="is-loading">
                                        <Loading />
                                    </el-icon> 处理中
                                </el-tag>
                            </template>
                        </el-table-column>
                    </el-table>
                </div>

                <div v-if="queueStatus.pending_embedding_tasks.length > 0" class="task-section">
                    <el-divider content-position="left" class="task-divider">
                        <el-tag type="info" size="small">等待中 ({{ queueStatus.pending_embedding_tasks.length }})</el-tag>
                    </el-divider>
                    <el-table :data="queueStatus.pending_embedding_tasks" size="small" stripe>
                        <el-table-column prop="opera_id" label="视频ID" width="100" />
                        <el-table-column prop="id" label="任务ID" show-overflow-tooltip />
                        <el-table-column prop="created_at" label="创建时间" width="180">
                            <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
                        </el-table-column>
                        <el-table-column label="状态" width="100">
                            <template #default>
                                <el-tag type="info" size="small">等待中</el-tag>
                            </template>
                        </el-table-column>
                    </el-table>
                </div>

                <div v-if="queueStatus.completed_embedding_tasks.length > 0" class="task-section">
                    <el-divider content-position="left" class="task-divider">
                        <el-tag type="primary" size="small">已完成 ({{ queueStatus.completed_embedding_tasks.length
                        }})</el-tag>
                    </el-divider>
                    <el-table :data="queueStatus.completed_embedding_tasks" size="small" stripe>
                        <el-table-column prop="opera_id" label="视频ID" width="100" />
                        <el-table-column prop="id" label="任务ID" show-overflow-tooltip />
                        <el-table-column prop="created_at" label="创建时间" width="180">
                            <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
                        </el-table-column>
                        <el-table-column prop="updated_at" label="完成时间" width="180">
                            <template #default="{ row }">{{ formatTime(row.updated_at) }}</template>
                        </el-table-column>
                        <el-table-column label="状态" width="100">
                            <template #default="{ row }">
                                <el-tag v-if="row.status === 'completed'" type="success" size="small">成功</el-tag>
                                <el-tag v-else-if="row.status === 'failed'" type="danger" size="small">失败</el-tag>
                            </template>
                        </el-table-column>
                        <el-table-column prop="error" label="错误信息" show-overflow-tooltip />
                    </el-table>
                </div>

                <el-empty
                    v-if="queueStatus.pending_embedding_tasks.length === 0 && queueStatus.processing_embedding_tasks.length === 0 && queueStatus.completed_embedding_tasks.length === 0"
                    description="暂无向量生成任务" :image-size="80" />
            </el-tab-pane>

            <el-tab-pane label="摘要生成任务">
                <template #label>
                    <span>
                        <el-icon>
                            <Document />
                        </el-icon>
                        摘要生成
                        <el-badge
                            v-if="queueStatus.pending_summary_tasks.length + queueStatus.processing_summary_tasks.length > 0"
                            :value="queueStatus.pending_summary_tasks.length + queueStatus.processing_summary_tasks.length"
                            class="item" />
                    </span>
                </template>

                <div v-if="queueStatus.processing_summary_tasks.length > 0" class="task-section">
                    <el-divider content-position="left" class="task-divider">
                        <el-tag type="success" size="small">正在处理 ({{ queueStatus.processing_summary_tasks.length
                        }})</el-tag>
                    </el-divider>
                    <el-table :data="queueStatus.processing_summary_tasks" size="small" stripe>
                        <el-table-column prop="opera_id" label="视频ID" width="100" />
                        <el-table-column prop="id" label="任务ID" show-overflow-tooltip />
                        <el-table-column prop="created_at" label="创建时间" width="180">
                            <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
                        </el-table-column>
                        <el-table-column prop="updated_at" label="更新时间" width="180">
                            <template #default="{ row }">{{ formatTime(row.updated_at) }}</template>
                        </el-table-column>
                        <el-table-column label="状态" width="100">
                            <template #default>
                                <el-tag type="success" size="small">
                                    <el-icon class="is-loading">
                                        <Loading />
                                    </el-icon> 处理中
                                </el-tag>
                            </template>
                        </el-table-column>
                    </el-table>
                </div>

                <div v-if="queueStatus.pending_summary_tasks.length > 0" class="task-section">
                    <el-divider content-position="left" class="task-divider">
                        <el-tag type="info" size="small">等待中 ({{ queueStatus.pending_summary_tasks.length }})</el-tag>
                    </el-divider>
                    <el-table :data="queueStatus.pending_summary_tasks" size="small" stripe>
                        <el-table-column prop="opera_id" label="视频ID" width="100" />
                        <el-table-column prop="id" label="任务ID" show-overflow-tooltip />
                        <el-table-column prop="created_at" label="创建时间" width="180">
                            <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
                        </el-table-column>
                        <el-table-column label="状态" width="100">
                            <template #default>
                                <el-tag type="info" size="small">等待中</el-tag>
                            </template>
                        </el-table-column>
                    </el-table>
                </div>

                <div v-if="queueStatus.completed_summary_tasks.length > 0" class="task-section">
                    <el-divider content-position="left" class="task-divider">
                        <el-tag type="primary" size="small">已完成 ({{ queueStatus.completed_summary_tasks.length
                        }})</el-tag>
                    </el-divider>
                    <el-table :data="queueStatus.completed_summary_tasks" size="small" stripe>
                        <el-table-column prop="opera_id" label="视频ID" width="100" />
                        <el-table-column prop="id" label="任务ID" show-overflow-tooltip />
                        <el-table-column prop="created_at" label="创建时间" width="180">
                            <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
                        </el-table-column>
                        <el-table-column prop="updated_at" label="完成时间" width="180">
                            <template #default="{ row }">{{ formatTime(row.updated_at) }}</template>
                        </el-table-column>
                        <el-table-column label="状态" width="100">
                            <template #default="{ row }">
                                <el-tag v-if="row.status === 'completed'" type="success" size="small">成功</el-tag>
                                <el-tag v-else-if="row.status === 'failed'" type="danger" size="small">失败</el-tag>
                            </template>
                        </el-table-column>
                        <el-table-column prop="error" label="错误信息" show-overflow-tooltip />
                    </el-table>
                </div>

                <el-empty
                    v-if="queueStatus.pending_summary_tasks.length === 0 && queueStatus.processing_summary_tasks.length === 0 && queueStatus.completed_summary_tasks.length === 0"
                    description="暂无摘要生成任务" :image-size="80" />
            </el-tab-pane>
        </el-tabs>
    </el-card>
</template>

<script lang="ts">
import TaskQueueViewScript from '../../scripts/views/admin/TaskQueueView';
export default TaskQueueViewScript;
</script>

<style scoped src="../../styles/views/admin/TaskQueueView.css"></style>