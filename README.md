# **Weak LINQ: Lazy LINQ-like Functions for Go**

An implementation of **lazy linq-like functions for Go**. This library provides a familiar, expressive, and declarative way to query and manipulate collections, 
inspired by C\#'s Language Integrated Query (LINQ). By utilizing **lazy evaluation**, it only processes elements when the result is materialized, 
improving performance for long chains of operations.

## **🚀 Features**

The weak-linq library aims to provide functionality similar to that of .NET's LINQ, but for Go, including:

* **Lazy Evaluation:** Operations are deferred until a materialization function is called (e.g., ToSlice()), allowing for efficient chaining of operations.  
* **Filtering:** Selectively includes elements from a collection (e.g., FilterOn).  
* **Transformation:** Projects elements into a new form (e.g., Get).  
* **Grouping:** Organizes elements into groups based on a key (e.g., GroupBy).  
* **Joining:** Combines elements from different collections (e.g., Join).  
* **Materialization:** Functions to convert the lazy query back into a concrete slice or map (e.g., AndAssignToSlice(), AndAssignToMap()).

weak-linq also has a few features not found in LINQ like:

* **Group Lists:** GroupListsBy functions allowing for quick grouping without overwriting values
* **Field Name Arg Option:** All functions have a version that accepts a field name for readability (expect a slight performance hit)
* **API Aiming For Human Readablity:** Functions are chanined together in a way that will hopefully make code more readable to others

## **📦 Installation**

To start using weak-linq, simply use go get to add the library to your project:
```
go get github.com/madicen/weak-linq/v2@latest
```

## **💡 Usage and Examples**

The entry point for most queries is the From() function, which initializes a new lazy query.

### **Human Readable Example Query**

The following example demonstrates grouping lists of people's names by their age.

```
type Person struct { Name string; Age int }
people := []Person{...}

peopleByAge := map[int][]string

weaklinq.From(people).
  GroupListsOf("Name").
  By("Age").
  AndAssignToMap(&peopleByAge)
```

### **Basic LINQ-style Example**

The following example demonstrates filtering out numbers less than or equal to 3 and then doubling the remaining numbers.

```
data := []int{1, 5, 2, 8, 3, 7}

// 1. Start the query from a slice.
// 2. Filter: Keep only elements greater than 3 (5, 8, 7).
// 3. Transform: Multiply each element by 2 (10, 16, 14).
// 4. Materialize: Convert the lazy query result back into a slice.

result := make([]int, 0)
weaklinq.From(data).
  FilterOnThis(func(x int) any {
    return x > 3
  }).
  GetThese(func(x int) int {
    return x * 2
  }).
  AndAssignToSlice(&result)

fmt.Println(result) 
// Output: [10 16 14]
```

## **🤝 Contributing**

Contributions are welcome\! If you would like to contribute, please:

1. Fork the repository.  
2. Create a new branch (git checkout \-b feature/awesome-thing).  
3. Commit your changes (git commit \-m 'Add awesome thing').  
4. Push to the branch (git push origin feature/awesome-thing).  
5. Open a Pull Request.

## **📄 License**

This project is licensed under the **MIT License**. See the [LICENSE](https://www.google.com/search?q=https://github.com/madicen/weak-linq/blob/main/LICENSE) file for details.
