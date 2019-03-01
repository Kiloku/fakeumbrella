import { Component, OnInit } from '@angular/core';
import { HttpClient,HttpHeaders  } from '@angular/common/http';
import {Observable} from "rxjs/Observable";
import { switchMap } from 'rxjs/operators';      
import { Router, ActivatedRoute, ParamMap } from '@angular/router';

@Component({
  selector: 'app-manage',
  templateUrl: './manage.component.html',
  styleUrls: ['./manage.component.css']
})
export class ManageComponent implements OnInit {
	customer$ : Observable<any>;
	customer : any;
	
	constructor(
		private http:HttpClient,
		private route: ActivatedRoute,
  		private router: Router){}

	ngOnInit(){
		let id = this.route.snapshot.paramMap.get('id');
		console.log(id);
		console.log(id);
		this.get_customer(id);
	}

	get_customer(id){
        this.http.get("http://localhost:8080/customers/" + id).subscribe((res : any)=>{
            console.log(res);
            this.customer = res;
        });
	}

	rename_customer(name){
		console.log("renaming to " + name);
		this.http.put("http://localhost:8080/customers/" + this.customer.id, {'name': name}).subscribe((res:any)=>{
			console.log(res);
			this.get_customer(this.customer.id);
		});
	}

	delete_customer(){
		this.http.delete("http://localhost:8080/customers/" + this.customer.id).subscribe((res:any)=>{
			console.log("deleting...");
			console.log(res);
			this.back();
		});
	}

	back(){
		this.router.navigate(['customers']);

	}

}
