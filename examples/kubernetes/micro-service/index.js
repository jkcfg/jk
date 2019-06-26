import * as param from '@jkcfg/std/param';
import { MicroService } from './micro-service';

const service = param.Object('service');
export default MicroService(service);
