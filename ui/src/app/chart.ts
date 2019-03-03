import {Engine, Scene, FreeCamera, Light,
Vector3, HemisphericLight, MeshBuilder, Mesh, 
StandardMaterial, Color3
} from 'babylonjs';
import {AdvancedDynamicTexture, TextBlock} from 'babylonjs-gui';
export class Chart 
{
	private _canvas: HTMLCanvasElement;
	private _engine: Engine;
	private _scene: Scene;
	private _camera: FreeCamera;
	private _light: Light;
	private _bars : Mesh[] = [];
	private _gui : AdvancedDynamicTexture;
	private _text : TextBlock[] = [];
	private _maxLabel : TextBlock;
	private _midLabel : TextBlock;
	private _minLabel : TextBlock;

	private _materialGreen : StandardMaterial;
	private _materialRed : StandardMaterial;

	constructor(canvasElement : string) {
		this._canvas = <HTMLCanvasElement>document.getElementById(canvasElement);
		this._engine = new Engine(this._canvas, true); 
	}

	createScene() : void {
		//this._bars = 
		this._scene = new Scene(this._engine);
		this._gui = AdvancedDynamicTexture.CreateFullscreenUI("myUI");
		this._camera = new FreeCamera('camera1', new Vector3(0,5,-10), this._scene);
		this._materialRed = new StandardMaterial('red', this._scene);
		this._materialGreen = new StandardMaterial('green', this._scene);
		this._materialRed.alpha = 1;
		this._materialGreen.alpha = 1;
		this._materialRed.diffuseColor = new Color3(0.8, 0.4, 0.4);
		this._materialGreen.diffuseColor = new Color3(0.4,  0.8, 0.4);

		this._camera.setTarget(Vector3.Zero());

        this._camera.mode = FreeCamera.ORTHOGRAPHIC_CAMERA;
        this._camera.orthoBottom = -1.5 * 4;
        this._camera.orthoTop = 1.5 * 4;
        this._camera.orthoLeft = -1.5 * 6;
        this._camera.orthoRight = 1.5 * 6;

		//this._camera.attachControl(this._canvas, false);

		this._light = new HemisphericLight('light1', new Vector3(0,1,0), this._scene);
		for (let i = 0; i < 4; i++)
		{
			this._bars[i] = MeshBuilder.CreateBox('bar'+i, {size:2}, this._scene);
			this._bars[i].position.x = -3.5 + (i * 3);
			this._bars[i].position.y = -3.5;
			this._bars[i].setPivotPoint(new Vector3(0,-1,0));
			this._bars[i].material = this._materialRed;
			this._text[i] = new TextBlock('text' + i, 'comp' + i);
			this._gui.addControl(this._text[i]);
			this._text[i].left = -2 + (this._bars[i].position.x * 33);
			this._text[i].top = 165;
		}
		this._maxLabel = new TextBlock('maxLabel', "###");
		this._midLabel = new TextBlock('maxLabel', "XXX");
		this._minLabel = new TextBlock('maxLabel', "0");

		this._maxLabel.left = -175;
		this._maxLabel.top = -150;

		this._midLabel.left = -175;
		this._midLabel.top = 0;

		this._minLabel.left = -175;
		this._minLabel.top = 150;

		this._gui.addControl(this._maxLabel);
		this._gui.addControl(this._midLabel);
		this._gui.addControl(this._minLabel);

		//let ground = MeshBuilder.CreateGround('ground1', {width: 30, height: 30, subdivisions: 2}, this._scene);
		let wall = MeshBuilder.CreateBox('wall', {size: 5}, this._scene);
		wall.position.z = 13;
		wall.scaling.x *= 10;
		wall.scaling.y *= 10;
		//ground.position.y = 0;
	}

	render(customers, forecasts) : void {
		this._engine.runRenderLoop(() => {
			this.pickCustomers(customers, forecasts);
    		this._scene.render();


			this._engine.stopRenderLoop()
		});

		window.addEventListener('resize', () => {
			this._engine.resize()
		});
	}

	pickCustomers(customers, forecasts) : void {
		customers = customers.sort((a, b) => {
			return b.employees - a.employees;
		});
		for (let i = 0; i < 4; i++ )
		{
			//console.log(customers[1].employees);
			//console.log(customers[i].name + "size: " + (customers[0].employees/customers[i].employees));

			this._bars[i].scaling.y = 5 * (customers[i].employees/customers[0].employees);
			this._text[i].text = (customers[i].name).split(' ').join("\n");

			this._maxLabel.text = customers[0].employees + ""; 
			this._midLabel.text = customers[0].employees/2 + "";




			for (let j = 0; j < forecasts.length; j++)
			{

				if (customers[i].location == forecasts[j].City && forecasts[j].Forecast.willRain)
				{
					this._bars[i].material = this._materialGreen;
				}
				/*console.log(forecasts[j].Forecast.willRain);

				{
					for (let h = 0; h < forecasts[j].Customers.length; h++)
					{
						console.log(forecasts[j].Customers[h].name + " - " + customers[i].name)
						if (forecasts[j].Customers[h].name == customers[i].name)
						{
							this._bars[i].material = this._materialRed;
						}
					}
				}*/
			}
		}
	}
}
