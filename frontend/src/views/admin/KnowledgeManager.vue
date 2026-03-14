<template>
  <div class="knowledge-manager">
    <!-- 统计卡片 -->
    <div class="stats-container">
      <el-row :gutter="16">
        <el-col :xs="12" :sm="8" :md="4">
          <div class="stat-card">
            <div class="stat-value">{{ stats.total_documents ?? 0 }}</div>
            <div class="stat-label">总文档数</div>
          </div>
        </el-col>
        <el-col :xs="12" :sm="8" :md="4">
          <div class="stat-card">
            <div class="stat-value">{{ stats.professional_docs ?? 0 }}</div>
            <div class="stat-label">专业文档</div>
          </div>
        </el-col>
        <el-col :xs="12" :sm="8" :md="4">
          <div class="stat-card">
            <div class="stat-value">{{ stats.opera_docs ?? 0 }}</div>
            <div class="stat-label">剧情文档</div>
          </div>
        </el-col>
        <el-col :xs="12" :sm="8" :md="4">
          <div class="stat-card stat-card--success">
            <div class="stat-value">{{ stats.completed_docs ?? 0 }}</div>
            <div class="stat-label">已向量化</div>
          </div>
        </el-col>
        <el-col :xs="12" :sm="8" :md="4">
          <div class="stat-card stat-card--warning">
            <div class="stat-value">{{ stats.pending_docs ?? 0 }}</div>
            <div class="stat-label">待处理</div>
          </div>
        </el-col>
        <el-col :xs="12" :sm="8" :md="4">
          <div class="stat-card stat-card--info">
            <div class="stat-value">{{ stats.total_chunks ?? 0 }}</div>
            <div class="stat-label">总分块数</div>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 工具栏 -->
    <div class="toolbar">
      <el-select v-model="filterSource" placeholder="来源筛选" style="width: 130px;" clearable @change="handleFilter">
        <el-option label="专业文档" value="professional" />
        <el-option label="剧情文档" value="opera" />
      </el-select>

      <el-select v-model="filterStatus" placeholder="状态筛选" style="width: 130px;" clearable @change="handleFilter">
        <el-option label="待处理" value="pending" />
        <el-option label="处理中" value="processing" />
        <el-option label="已完成" value="completed" />
        <el-option label="失败" value="failed" />
      </el-select>

      <div class="toolbar-right">
        <el-button :icon="Refresh" circle @click="handleRefresh" :loading="loading" title="刷新列表" />
        <el-button @click="handleImportOperas" :loading="importing" type="default">
          导入作品字幕
        </el-button>
        <el-button type="primary" @click="handleShowUploadDialog">
          上传知识文档
        </el-button>
      </div>
    </div>

    <!-- 文档表格 -->
    <el-table :data="filteredData" style="width: 100%" v-loading="loading" stripe>
      <el-table-column prop="title" label="文档标题" min-width="220" show-overflow-tooltip />

      <el-table-column label="来源" width="90" align="center">
        <template #default="scope">
          <el-tag :type="scope.row.source_type === 'opera' ? 'primary' : 'success'" size="small">
            {{ scope.row.source_type === 'opera' ? '剧情' : '专业' }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column label="向量化状态" width="115" align="center">
        <template #default="scope">
          <el-tag :type="getStatusType(scope.row.embedding_status)" effect="plain" size="small">
            {{ formatStatus(scope.row.embedding_status) }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column prop="chunks_count" label="分块数" width="80" align="center" />

      <el-table-column prop="created_at" label="创建时间" width="165">
        <template #default="scope">
          {{ formatTime(scope.row.created_at) }}
        </template>
      </el-table-column>

      <el-table-column label="操作" width="200" fixed="right">
        <template #default="scope">
          <el-button size="small" @click="handleViewDetail(scope.row)">查看</el-button>
          <el-popconfirm
            title="确认删除该文档？删除后可在下方重新激活。"
            confirm-button-text="删除"
            cancel-button-text="取消"
            @confirm="handleDeleteDocument(scope.row.doc_id)"
          >
            <template #reference>
              <el-button size="small" type="danger">删除</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <el-pagination
      style="margin-top: 20px; justify-content: flex-end;"
      :current-page="currentPage"
      :page-size="pageSize"
      :page-sizes="[10, 20, 50]"
      :total="total"
      layout="total, sizes, prev, pager, next"
      @current-change="handlePageChange"
      @size-change="handlePageSizeChange"
    />

    <!-- 上传文档对话框 -->
    <el-dialog v-model="uploadDialogVisible" title="上传知识文档" width="55%" :close-on-click-modal="false">
      <el-form
        ref="uploadFormRef"
        :model="uploadForm"
        label-width="90px"
        :rules="uploadFormRules"
      >
        <el-form-item label="文档标题" prop="title">
          <el-input v-model="uploadForm.title" placeholder="请输入文档标题" maxlength="255" show-word-limit />
        </el-form-item>

        <el-form-item label="文档来源" prop="source_type">
          <el-radio-group v-model="uploadForm.source_type">
            <el-radio value="professional">专业文档</el-radio>
            <el-radio value="opera">剧情文档</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="文档内容" prop="content">
          <el-input
            v-model="uploadForm.content"
            type="textarea"
            :rows="12"
            placeholder="粘贴或输入文档内容（支持纯文本和 Markdown）"
            maxlength="100000"
            show-word-limit
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="uploadDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitUpload" :loading="uploading">上传</el-button>
      </template>
    </el-dialog>

    <!-- 文档详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="文档详情" width="70%" :close-on-click-modal="false">
      <div v-if="currentDocument" class="document-detail">
        <!-- 基本信息 -->
        <el-descriptions :column="2" border size="small">
          <el-descriptions-item label="文档标题" :span="2">{{ currentDocument.title }}</el-descriptions-item>
          <el-descriptions-item label="文档来源">
            <el-tag :type="currentDocument.source_type === 'opera' ? 'primary' : 'success'" size="small">
              {{ currentDocument.source_type === 'opera' ? '剧情文档' : '专业文档' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="向量化状态">
            <el-tag :type="getStatusType(currentDocument.embedding_status)" size="small">
              {{ formatStatus(currentDocument.embedding_status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="分块数">{{ currentDocument.chunks_count }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ formatTime(currentDocument.created_at) }}</el-descriptions-item>
          <el-descriptions-item label="文档 ID" :span="2">
            <span class="doc-id">{{ currentDocument.doc_id }}</span>
          </el-descriptions-item>
        </el-descriptions>

        <!-- 内容 Tab -->
        <el-tabs v-model="detailActiveTab" class="content-tabs">
          <!-- 完整原文 -->
          <el-tab-pane label="完整原文" name="full">
            <div class="tab-toolbar">
              <span class="content-length">共 {{ currentDocument.content?.length ?? 0 }} 字符</span>
            </div>
            <div class="markdown-body" v-html="fullContentHtml" />
          </el-tab-pane>

          <!-- 分块列表 -->
          <el-tab-pane :label="`向量分块（${currentDocument.chunks?.length ?? 0}）`" name="chunks">
            <div v-if="!currentDocument.chunks?.length" class="no-chunks">
              暂无分块数据（文档尚未完成向量化）
            </div>
            <el-collapse v-else accordion class="chunk-collapse">
              <el-collapse-item
                v-for="(chunk, idx) in currentDocument.chunks"
                :key="chunk.chunk_id"
                :name="idx"
              >
                <template #title>
                  <span class="chunk-title">块 {{ chunk.chunk_index + 1 }}</span>
                  <span class="chunk-length">{{ chunk.chunk_text.length }} 字符</span>
                </template>
                <div class="markdown-body" v-html="chunksHtml[idx]" />
              </el-collapse-item>
            </el-collapse>
          </el-tab-pane>
        </el-tabs>
      </div>
    </el-dialog>
  </div>
</template>

<script lang="ts" src="../../scripts/views/admin/KnowledgeManager.ts"></script>
<style scoped src="../../styles/views/admin/KnowledgeManager.css"></style>
