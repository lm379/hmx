<template>
  <div class="education-manager">
    <div class="toolbar">
      <el-button type="primary" @click="handleCreate">上传书籍 PDF</el-button>
      <el-button @click="fetchData">刷新</el-button>
    </div>

    <el-table :data="tableData" style="width: 100%" v-loading="loading">
      <el-table-column prop="book_id" label="ID" width="80" />
      <el-table-column label="封面" width="110">
        <template #default="scope">
          <el-image
            v-if="scope.row.cover_url"
            :src="scope.row.cover_url"
            style="width: 64px; height: 88px; border-radius: 4px;"
            fit="cover"
          />
          <span v-else class="muted-text">无封面</span>
        </template>
      </el-table-column>
      <el-table-column prop="title" label="书籍名称" min-width="220" show-overflow-tooltip />
      <el-table-column prop="description" label="简介" min-width="260" show-overflow-tooltip />
      <el-table-column label="状态" width="100">
        <template #default="scope">
          <el-tag :type="scope.row.is_published ? 'success' : 'info'">
            {{ scope.row.is_published ? '已发布' : '草稿' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="180">
        <template #default="scope">{{ formatTime(scope.row.created_at) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="260" fixed="right">
        <template #default="scope">
          <el-button size="small" @click="openPdf(scope.row)">预览</el-button>
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

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑书籍' : '上传书籍 PDF'" width="52%" @paste="handlePaste">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="90px">
        <el-form-item label="书籍名称" prop="title">
          <el-input v-model="form.title" maxlength="200" show-word-limit />
        </el-form-item>
        <el-form-item label="简介">
          <el-input v-model="form.description" type="textarea" :rows="4" />
        </el-form-item>
        <el-form-item label="PDF 文件" prop="pdf_path">
          <div class="pdf-upload-field">
            <div v-if="form.pdf_path && !pdfFile" class="current-pdf">
              当前文件已上传，可直接保存或重新选择 PDF
            </div>
            <el-upload
              ref="pdfUploadRef"
              action="#"
              :auto-upload="false"
              :limit="1"
              :on-change="handlePdfChange"
              :on-remove="handlePdfRemove"
              :on-exceed="handlePdfExceed"
              accept="application/pdf,.pdf"
              drag
            >
              <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
              <div class="el-upload__text">
                将 PDF 拖到此处，或<em>点击上传</em>
              </div>
              <template #tip>
                <div class="el-upload__tip">只支持 PDF，保存时自动上传</div>
              </template>
            </el-upload>
          </div>
        </el-form-item>
        <el-form-item label="PDF 封面">
          <div class="cover-edit-field">
            <div class="cover-preview-field">
              <img v-if="coverPreviewUrl" :src="coverPreviewUrl" class="cover-preview-img" alt="PDF 封面预览" />
              <div v-else class="cover-placeholder">
                {{ coverGenerating ? '正在生成封面...' : '选择 PDF 后自动生成第一页封面' }}
              </div>
            </div>
            <el-upload
              ref="coverUploadRef"
              action="#"
              :auto-upload="false"
              :limit="1"
              :on-change="handleCoverChange"
              :on-remove="handleCoverRemove"
              :on-exceed="handleCoverExceed"
              accept="image/*"
            >
              <el-button>自定义封面</el-button>
              <template #tip>
                <div class="el-upload__tip">
                  支持 JPG、PNG 等图片，也可在对话框内粘贴图片；不上传则默认使用 PDF 第一页
                </div>
              </template>
            </el-upload>
            <el-tag v-if="coverMode === 'custom'" type="success" effect="plain">自定义封面</el-tag>
            <el-tag v-else-if="coverMode === 'auto'" type="info" effect="plain">PDF 首页封面</el-tag>
          </div>
        </el-form-item>
        <el-form-item label="发布状态">
          <el-switch v-model="form.is_published" active-text="发布" inactive-text="草稿" />
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
import EducationManagerScript from '../../scripts/views/admin/EducationManager';
import { UploadFilled } from '@element-plus/icons-vue';

export default {
  ...EducationManagerScript,
  components: {
    UploadFilled
  }
};
</script>

<style scoped src="../../styles/views/admin/EducationManager.css"></style>
