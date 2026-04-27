<template>
  <div class="news-manager">
    <div class="toolbar">
      <el-button type="primary" @click="handleCreate">新增新闻</el-button>
      <el-button @click="fetchData">刷新</el-button>
    </div>

    <el-table :data="tableData" style="width: 100%" v-loading="loading">
      <el-table-column prop="news_id" label="ID" width="80" />
      <el-table-column label="封面" width="120">
        <template #default="scope">
          <el-image
            v-if="scope.row.cover"
            :src="scope.row.cover"
            style="width: 80px; height: 45px; border-radius: 4px;"
            fit="cover"
          />
          <span v-else class="muted-text">无封面</span>
        </template>
      </el-table-column>
      <el-table-column prop="title" label="标题" min-width="220" show-overflow-tooltip />
      <el-table-column prop="source" label="来源" width="120" />
      <el-table-column label="状态" width="100">
        <template #default="scope">
          <el-tag :type="scope.row.is_published ? 'success' : 'info'">
            {{ scope.row.is_published ? '已发布' : '草稿' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="发布时间" width="180">
        <template #default="scope">
          {{ formatTime(scope.row.published_at || scope.row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="220" fixed="right">
        <template #default="scope">
          <el-button size="small" @click="handleEdit(scope.row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        layout="prev, pager, next"
        @current-change="fetchData"
      />
    </div>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑新闻' : '新增新闻'" width="62%" @paste="handlePaste">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="90px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" maxlength="200" show-word-limit />
        </el-form-item>
        <el-form-item label="摘要" prop="summary">
          <el-input v-model="form.summary" type="textarea" :rows="3" maxlength="500" show-word-limit />
        </el-form-item>
        <el-form-item label="正文" prop="content">
          <el-input
            v-model="form.content"
            type="textarea"
            :rows="12"
            placeholder="支持 Markdown，也可以直接输入普通文本"
          />
        </el-form-item>
        <el-form-item label="封面图片">
          <div class="cover-upload-field">
            <el-image
              v-if="form.cover && !coverFile"
              :src="form.cover"
              class="cover-preview"
              fit="cover"
            />
            <el-upload
              ref="coverUploadRef"
              action="#"
              :auto-upload="false"
              :limit="1"
              :on-change="handleCoverChange"
              :on-remove="handleCoverRemove"
              :on-exceed="handleCoverExceed"
              list-type="picture"
              drag
            >
              <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
              <div class="el-upload__text">
                将封面图片拖到此处，或<em>点击上传</em>
              </div>
              <template #tip>
                <div class="el-upload__tip">支持拖拽、点击上传，或在对话框内直接粘贴剪贴板图片</div>
              </template>
            </el-upload>
          </div>
        </el-form-item>
        <el-form-item label="来源">
          <el-input v-model="form.source" placeholder="例如：平台资讯、文旅活动、媒体报道" />
        </el-form-item>
        <el-form-item label="作者">
          <el-input v-model="form.author" />
        </el-form-item>
        <el-form-item label="发布状态">
          <el-switch v-model="form.is_published" active-text="发布" inactive-text="草稿" />
        </el-form-item>
        <el-form-item label="发布时间">
          <el-date-picker
            v-model="form.published_at"
            type="datetime"
            value-format="YYYY-MM-DDTHH:mm:ssZ"
            placeholder="不填则发布时使用当前时间"
            style="width: 100%;"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" :loading="submitting" @click="handleSubmit">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import NewsManagerScript from '../../scripts/views/admin/NewsManager';
import { UploadFilled } from '@element-plus/icons-vue';

export default {
  ...NewsManagerScript,
  components: {
    UploadFilled
  }
};
</script>

<style scoped src="../../styles/views/admin/NewsManager.css"></style>
