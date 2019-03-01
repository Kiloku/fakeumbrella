import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import {Observable} from "rxjs/Observable";
@Component({
  selector: 'app-customers',
  templateUrl: './customers.component.html',
  styleUrls: ['./customers.component.css']
})
export class CustomersComponent implements OnInit {
	title = "Fake Umbrella";
	customers : Observable<any[]>;
	constructor(private http:HttpClient){}
	ngOnInit(){
		this.customers = this.http.get<any[]>("http://localhost:8080/customers/")
	}


}
