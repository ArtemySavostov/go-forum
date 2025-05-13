//ParentComponent.js (where you're using ProductList)
import ProductList from './ProductList';

const ParentComponent = () => {
  const handleCreateProduct = (productData) => {
    // Do something with the product data (e.g., send it to an API)
    console.log('Creating product:', productData);
  };

  return (
    <div>
      <ProductList onCreate={handleCreateProduct} />
    </div>
  );
};

export default ParentComponent;

// import ProductList from './ProductList';

// const ParentComponent = ({ products, onCreate }) => {  // Receive props
//   const handleCreateProduct = (productData) => {
//     onCreate(productData); // Pass the data to the App component's onCreate
//   };

//   return (
//     <div>
//       <ProductList products={products} onCreate={handleCreateProduct} />  // Pass props down to ProductList
//     </div>
//   );
// };

// export default ParentComponent;
