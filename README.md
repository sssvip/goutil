# goutil
> simple go util

### testutil
preview:

![h](testutil/test.svg)

example:

```golang

func testStepA() error {
	//balabala
	//mock pass test
	return nil // pass step A
}

var stepBTryTimes = 0

func testStepB() error {
	stepBTryTimes++
	//step b is a unstable flow,need try to pass test
	if stepBTryTimes >= 5 {
		return nil //mock unstable situation
	}
	return errors.New("has error") // step b has error
}

func testStepC() error {
	//mock test failed
	return errors.New("has error") // step b has error
}

func main() {
	testutil.StartTest()
	//just show test
	testutil.TryMoreTime(testStepA, 1, "testStepA")  //should pass
	testutil.TryMoreTime(testStepB, 10, "testStepB") // try 10 times to pass step b, should pass.
	testutil.TryMoreTime(testStepC, 3, "testStepC")  // try 3 times to pass step c, should failed.
	testutil.EndTest()
}


```


### sql
1. sql generator
    
2. db helper

### http

### string

### log