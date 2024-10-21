package simple

import "html/template"

var repositoryTmpl = template.Must(template.New("repository").Parse(`
package repositories

import (
	"{{.PkgName}}/model"
	"github.com/jimersylee/go-bbs/utils/simple"
	"github.com/jinzhu/gorm"
)

var {{.Name}}Repository = new{{.Name}}Repository()

func new{{.Name}}Repository() *{{.CamelName}}Repository {
	return &{{.CamelName}}Repository{}
}

type {{.CamelName}}Repository struct {
}

func (this *{{.CamelName}}Repository) Get(db *gorm.DB, id int64) *model.{{.Name}} {
	ret := &model.{{.Name}}{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *{{.CamelName}}Repository) Take(db *gorm.DB, where ...interface{}) *model.{{.Name}} {
	ret := &model.{{.Name}}{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *{{.CamelName}}Repository) QueryCnd(db *gorm.DB, cnd *simple.QueryCnd) (list []model.{{.Name}}, err error) {
	err = cnd.DoQuery(db).Find(&list).Error
	return
}

func (this *{{.CamelName}}Repository) Query(db *gorm.DB, queries *simple.ParamQueries) (list []model.{{.Name}}, paging *simple.Paging) {
	queries.StartQuery(db).Find(&list)
    queries.StartCount(db).Model(&model.{{.Name}}{}).Count(&queries.Paging.Total)
	paging = queries.Paging
	return
}

func (this *{{.CamelName}}Repository) Create(db *gorm.DB, t *model.{{.Name}}) (err error) {
	err = db.Create(t).Error
	return
}

func (this *{{.CamelName}}Repository) Update(db *gorm.DB, t *model.{{.Name}}) (err error) {
	err = db.Save(t).Error
	return
}

func (this *{{.CamelName}}Repository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.{{.Name}}{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *{{.CamelName}}Repository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.{{.Name}}{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *{{.CamelName}}Repository) Delete(db *gorm.DB, id int64) {
	db.Model(&model.{{.Name}}{}).Delete("id", id)
}

`))

var serviceTmpl = template.Must(template.New("service").Parse(`
package services

import (
	"{{.PkgName}}/model"
	"{{.PkgName}}/repositories"
	"github.com/jimersylee/go-bbs/utils/simple"
)

var {{.Name}}Service = new{{.Name}}Service()

func new{{.Name}}Service() *{{.CamelName}}Service {
	return &{{.CamelName}}Service {}
}

type {{.CamelName}}Service struct {
}

func (this *{{.CamelName}}Service) Get(id int64) *model.{{.Name}} {
	return repositories.{{.Name}}Repository.Get(simple.GetDB(), id)
}

func (this *{{.CamelName}}Service) Take(where ...interface{}) *model.{{.Name}} {
	return repositories.{{.Name}}Repository.Take(simple.GetDB(), where...)
}

func (this *{{.CamelName}}Service) QueryCnd(cnd *simple.QueryCnd) (list []model.{{.Name}}, err error) {
	return repositories.{{.Name}}Repository.QueryCnd(simple.GetDB(), cnd)
}

func (this *{{.CamelName}}Service) Query(queries *simple.ParamQueries) (list []model.{{.Name}}, paging *simple.Paging) {
	return repositories.{{.Name}}Repository.Query(simple.GetDB(), queries)
}

func (this *{{.CamelName}}Service) Create(t *model.{{.Name}}) error {
	return repositories.{{.Name}}Repository.Create(simple.GetDB(), t)
}

func (this *{{.CamelName}}Service) Update(t *model.{{.Name}}) error {
	return repositories.{{.Name}}Repository.Update(simple.GetDB(), t)
}

func (this *{{.CamelName}}Service) Updates(id int64, columns map[string]interface{}) error {
	return repositories.{{.Name}}Repository.Updates(simple.GetDB(), id, columns)
}

func (this *{{.CamelName}}Service) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.{{.Name}}Repository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *{{.CamelName}}Service) Delete(id int64) {
	repositories.{{.Name}}Repository.Delete(simple.GetDB(), id)
}

`))

var controllerTmpl = template.Must(template.New("controller").Parse(`
package admin

import (
	"{{.PkgName}}/model"
	"{{.PkgName}}/services"
	"github.com/jimersylee/go-bbs/utils/simple"
	"github.com/kataras/iris/v12"
	"strconv"
)

type {{.Name}}Controller struct {
	Ctx             *context.Context
}

func (this *{{.Name}}Controller) GetBy(id int64) *simple.JsonResult {
	t := services.{{.Name}}Service.Get(id)
	if t == nil {
		return simple.ErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *{{.Name}}Controller) AnyList() *simple.JsonResult {
	list, paging := services.{{.Name}}Service.Query(simple.NewParamQueries(this.Ctx).PageAuto().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (this *{{.Name}}Controller) PostCreate() *simple.JsonResult {
	t := &model.{{.Name}}{}
	err := this.Ctx.ReadForm(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}

	err = services.{{.Name}}Service.Create(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (this *{{.Name}}Controller) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	t := services.{{.Name}}Service.Get(id)
	if t == nil {
		return simple.ErrorMsg("entity not found")
	}

	err = this.Ctx.ReadForm(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}

	err = services.{{.Name}}Service.Update(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

`))

var viewIndexTmpl = template.Must(template.New("index.vue").Parse(`
<template>
    <section>
        <!--工具条-->
        <el-col :span="24" class="toolbar" style="padding-bottom: 0px;">
            <el-form :inline="true" :model="filters">
                <el-form-item>
                    <el-input v-model="filters.name" placeholder="名称"></el-input>
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" v-on:click="list">查询</el-button>
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" @click="handleAdd">新增</el-button>
                </el-form-item>
            </el-form>
        </el-col>

        <!--列表-->
        <el-table :data="results" highlight-current-row border v-loading="listLoading"
                  style="width: 100%;" @selection-change="handleSelectionChange">
            <el-table-column type="selection" width="55"></el-table-column>
            <el-table-column prop="id" label="编号"></el-table-column>
            {{range .Fields}}
                <el-table-column prop="{{.CamelName}}" label="{{.CamelName}}"></el-table-column>
            {{end}}
            <el-table-column label="操作" width="150">
                <template scope="scope">
                    <el-button size="small" @click="handleEdit(scope.$index, scope.row)">编辑</el-button>
                </template>
            </el-table-column>
        </el-table>

        <!--工具条-->
        <el-col :span="24" class="toolbar">
            <el-pagination layout="total, sizes, prev, pager, next, jumper" :page-sizes="[20, 50, 100, 300]"
                           @current-change="handlePageChange"
                           @size-change="handleLimitChange"
                           :current-page="page.page"
                           :page-size="page.limit"
                           :total="page.total"
                           style="float:right;">
            </el-pagination>
        </el-col>


        <!--新增界面-->
        <el-dialog title="新增" :visible.sync="addFormVisible" :close-on-click-modal="false">
            <el-form :model="addForm" label-width="80px" :rules="addFormRules" ref="addForm">
                {{range .Fields}}
                    <el-form-item label="{{.CamelName}}" prop="rule">
                        <el-input v-model="addForm.{{.CamelName}}"></el-input>
                    </el-form-item>
                {{end}}
            </el-form>
            <div slot="footer" class="dialog-footer">
                <el-button @click.native="addFormVisible = false">取消</el-button>
                <el-button type="primary" @click.native="addSubmit" :loading="addLoading">提交</el-button>
            </div>
        </el-dialog>

        <!--编辑界面-->
        <el-dialog title="编辑" :visible.sync="editFormVisible" :close-on-click-modal="false">
            <el-form :model="editForm" label-width="80px" :rules="editFormRules" ref="editForm">
                <el-input v-model="editForm.id" type="hidden"></el-input>
                {{range .Fields}}
                    <el-form-item label="{{.CamelName}}" prop="rule">
                        <el-input v-model="editForm.{{.CamelName}}"></el-input>
                    </el-form-item>
                {{end}}
            </el-form>
            <div slot="footer" class="dialog-footer">
                <el-button @click.native="editFormVisible = false">取消</el-button>
                <el-button type="primary" @click.native="editSubmit" :loading="editLoading">提交</el-button>
            </div>
        </el-dialog>
    </section>
</template>

<script>
  import HttpClient from '../../apis/HttpClient'

  export default {
    name: "List",
    data() {
      return {
        results: [],
        listLoading: false,
        page: {},
        filters: {},
        selectedRows: [],

        addForm: {
          {{range .Fields}}
            '{{.CamelName}}': '',
          {{end}}
        },
        addFormVisible: false,
        addFormRules: {},
        addLoading: false,

        editForm: {
          'id': '',
          {{range .Fields}}
            '{{.CamelName}}': '',
          {{end}}
        },
        editFormVisible: false,
        editFormRules: {},
        editLoading: false,
      }
    },
    mounted() {
      this.list();
    },
    methods: {
      list() {
        let me = this
        me.listLoading = true
		let params = Object.assign(me.filters, {
          page: me.page.page,
          limit: me.page.limit
        })
        HttpClient.post('/api/admin/{{.KebabName}}/list', params)
          .then(data => {
            me.results = data.results
            me.page = data.page
          })
          .finally(() => {
            me.listLoading = false
          })
      },
      handlePageChange (val) {
        this.page.page = val
        this.list()
      },
      handleLimitChange (val) {
        this.page.limit = val
        this.list()
      },
      handleAdd() {
        this.addForm = {
          name: '',
          description: '',
        }
        this.addFormVisible = true
      },
      addSubmit() {
        let me = this
        HttpClient.post('/api/admin/{{.KebabName}}/create', this.addForm)
          .then(data => {
            me.$message({message: '提交成功', type: 'success'});
            me.addFormVisible = false
            me.list()
          })
          .catch(rsp => {
            me.$notify.error({title: '错误', message: rsp.message})
          })
      },
      handleEdit(index, row) {
        let me = this
        HttpClient.get('/api/admin/{{.KebabName}}/' + row.id)
          .then(data => {
            me.editForm = Object.assign({}, data);
            me.editFormVisible = true
          })
          .catch(rsp => {
            me.$notify.error({title: '错误', message: rsp.message})
          })
      },
      editSubmit() {
        let me = this
        HttpClient.post('/api/admin/{{.KebabName}}/update', me.editForm)
          .then(data => {
            me.list()
            me.editFormVisible = false
          })
          .catch(rsp => {
            me.$notify.error({title: '错误', message: rsp.message})
          })
      },

      handleSelectionChange(val) {
        this.selectedRows = val
      },
    }
  }
</script>

<style scoped>

</style>

`))
