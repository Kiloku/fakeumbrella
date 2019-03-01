import { Component, OnInit } from '@angular/core';
import { HttpClient,HttpHeaders  } from '@angular/common/http';
import {Observable} from "rxjs/Observable";
import { switchMap } from 'rxjs/operators';      
import { Router, ActivatedRoute, ParamMap } from '@angular/router';

@Component({
  selector: 'app-new',
  templateUrl: './new.component.html',
  styleUrls: ['./new.component.css']
})
export class NewComponent implements OnInit {
	customer$ : Observable<any>;
	customer : any;

	constructor(
		private http: HttpClient,
		private route: ActivatedRoute,
  		private router: Router) { }

	ngOnInit() {
	}

	create_customer(name, location, employees, contact, telephone){
		console.log("creating " + name);
		this.http.post("http://localhost:8080/customers", 
			{'name': name, 'location':location, 'employees': parseInt(employees), 'contact':contact, 'telephone':telephone}, ).subscribe(
				(res:any)=>
			{
				this.router.navigate(['customers'])
				//console.log(res);

			});
	}
}
