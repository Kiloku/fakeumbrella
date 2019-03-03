import { Component, OnInit, ElementRef, AfterViewInit, ViewChild } from '@angular/core';
import * as BABYLON from 'babylonjs';
import { HttpClient } from '@angular/common/http';
import {Observable} from "rxjs/Observable"; 
import { Chart } from '../chart';
@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {
	//forecasts$ : Observable<any[]>;
	forecasts : object[];
	customers : object[];
	chart : Chart;
	constructor(private http:HttpClient){}

	ngOnInit(){
		this.getForecasts();
		this.http.get("http://localhost:8080/customers/").subscribe((res : any) => {
			this.customers = res;
			console.log(this.customers);
		});
	}

	ngAfterViewInit() {
		this.chart = new Chart('renderCanvas');
		this.chart.createScene();
	}

	getForecasts(){
        this.http.get("http://localhost:8080/forecasts/customers/").subscribe((res : any)=>{
            this.forecasts = res;
            /*this.forecasts.sort((a, b) => {
            	if(a.Forecast.willRain && b.Forecast.willRain) return 0;
            	if(a.Forecast.willRain && !b.Forecast.willRain) return 1;
            	if(!a.Forecast.willRain && b.Forecast.willRain) return 1;
            })
            this.forecasts.reverse();*/
            for (let i = 0; i < this.forecasts.length; i++)
            {
            	if (!this.forecasts[i]["Forecast"].willRain)
            	{
            		console.log("splicing");
            		let temp = this.forecasts[i];
            		this.forecasts.splice(i, 1)
            		this.forecasts.push(temp);
            	}
            }
            console.log(this.forecasts);
			this.chart.render(this.customers, this.forecasts);
        });

	}

	parseTime(time){
		let d = new Date(0);
		d.setUTCSeconds(time);

		return d;
	}


}
