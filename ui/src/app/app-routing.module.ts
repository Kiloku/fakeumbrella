import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import {HomeComponent} from './home/home.component';
import {CustomersComponent} from './customers/customers.component';
import {ManageComponent} from './customers/manage/manage.component';
import {NewComponent} from './customers/manage/new/new.component';

const routes: Routes = [
	{
		path: '',
		component: HomeComponent
	},
	{
		path: 'customers',
		component: CustomersComponent
	},
	{
		path: 'customers/manage/:id',
		component: ManageComponent,

	},
	{
		path: 'customers/new',
		component: NewComponent,
	},
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
